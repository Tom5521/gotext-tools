package parse

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Tom5521/gotext-tools/v2/internal/slices"
	"github.com/Tom5521/gotext-tools/v2/internal/util"
	"github.com/Tom5521/gotext-tools/v2/pkg/po"

	"github.com/alecthomas/participle/v2/lexer"
)

// Ensure PoParser implements the po.Parser interface.
var _ po.Parser = (*PoParser)(nil)

// PoParser handles parsing of PO files into po.File structures.
type PoParser struct {
	Config PoConfig // Configuration for parsing behavior

	originalData []byte // Original unmodified file data
	data         []byte // Current working data (may be modified during parsing)
	filename     string // Name of the source file

	errors []error // Collection of errors encountered during parsing
	warns  []error // Collection of non-critical warnings
}

// NewPo creates a new PoParser from a file path.
func NewPo(path string, options ...PoOption) (*PoParser, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return NewPoFromBytes(file, path, options...), nil
}

// NewPoFromReader creates a new PoParser from an io.Reader.
func NewPoFromReader(r io.Reader, name string, options ...PoOption) (*PoParser, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("error reading all reader: %w", err)
	}

	return NewPoFromBytes(data, name, options...), nil
}

// NewPoFromFile creates a new PoParser from an open *os.File.
func NewPoFromFile(f *os.File, options ...PoOption) (*PoParser, error) {
	return NewPoFromReader(f, f.Name(), options...)
}

// NewPoFromString creates a new PoParser from a string.
func NewPoFromString(s, name string, options ...PoOption) *PoParser {
	return NewPoFromBytes([]byte(s), name, options...)
}

// NewPoFromBytes creates a new PoParser from a byte slice.
func NewPoFromBytes(data []byte, name string, options ...PoOption) *PoParser {
	return &PoParser{
		Config:       DefaultPoConfig(options...),
		originalData: data,
		filename:     name,
	}
}

// error logs an error message and adds it to the parser's error collection.
// If a logger is configured, it will also log the error.
func (p *PoParser) error(format string, a ...any) {
	var err error
	format = "parse: " + format
	if len(a) == 0 {
		err = errors.New(format)
	} else {
		err = fmt.Errorf(format, a...)
	}

	if p.Config.Logger != nil {
		p.Config.Logger.Println("ERROR: ", err)
	}

	p.errors = append(p.errors, err)
}

// warn logs a warning message and adds it to the parser's warning collection.
// If verbose logging is enabled, it will also log the warning.
func (p *PoParser) warn(format string, a ...any) {
	var err error
	format = "warning: " + format
	if len(a) == 0 {
		err = errors.New(format)
	} else {
		err = fmt.Errorf(format, a...)
	}

	if p.Config.Logger != nil && p.Config.Verbose {
		p.Config.Logger.Println(err)
	}

	p.errors = append(p.errors, err)
}

// Error returns the first error encountered during parsing, if any.
func (p *PoParser) Error() error {
	if len(p.errors) == 0 {
		return nil
	}
	return p.errors[0]
}

// Warnings returns all non-critical warnings encountered during parsing.
func (p *PoParser) Warnings() []error {
	return p.warns
}

// Errors returns all errors encountered during parsing.
func (p *PoParser) Errors() []error {
	return p.errors
}

// Regular expressions for parsing different types of PO file comments.
var (
	locationRegex  = regexp.MustCompile(`#: *(.*)`)  // Source location comments
	generalRegex   = regexp.MustCompile(`# *(.*)`)   // General translator comments
	extractedRegex = regexp.MustCompile(`#\. *(.*)`) // Extracted comments
	flagRegex      = regexp.MustCompile(`#, *(.*)`)  // Flag comments
	obsoleteRegex  = regexp.MustCompile(`#~ *(.*)`)  // Obsolete entry markers
	previousRegex  = regexp.MustCompile(`#\| *(.*)`) // Previous message comments
)

// parseObsoleteEntries extracts and parses obsolete entries from the PO file.
func (p *PoParser) parseObsoleteEntries(tokens []lexer.Token) po.Entries {
	comp := obsoleteRegex
	if p.Config.UseCustomObsoletePrefix {
		var err error
		comp, err = regexp.Compile(fmt.Sprintf(`#%s *(.*)`, string(p.Config.CustomObsoletePrefix)))
		if err != nil {
			p.warn(err.Error())
			return nil
		}
	}

	cleanedLines := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if token.Type != util.PoSymbols["Comment"] {
			continue
		}
		if !comp.MatchString(token.String()) {
			continue
		}

		cleanedLines = append(cleanedLines,
			comp.FindStringSubmatch(token.String())[1],
		)
	}

	fullFile := strings.Join(cleanedLines, "\n")

	parser := NewPoFromString(
		fullFile,
		"tmp",
		PoWithConfig(p.Config),
		poWithMarkAllAsObsolete(true),
	)

	f := parser.Parse()
	err := parser.Error()
	if err != nil {
		return nil
	}

	return f.Entries
}

// parseComments processes all comment tokens associated with an entry and
// populates the corresponding fields in the Entry struct.
func parseComments(entry *po.Entry, tokens []lexer.Token) (err error) {
	for _, t := range tokens {
		if t.Type != util.PoSymbols["Comment"] {
			continue
		}
		if obsoleteRegex.MatchString(t.String()) {
			continue
		}
		switch {
		case locationRegex.MatchString(t.String()):
			matches := locationRegex.FindStringSubmatch(t.String())
			parts := strings.SplitN(matches[1], ":", 2)
			line := -1
			if parts[1] != "" {
				line, err = strconv.Atoi(parts[1])
				if err != nil {
					return fmt.Errorf("error parsing integer: %w", err)
				}
			}

			loc := po.Location{
				Line: line,
				File: parts[0],
			}
			entry.Locations = append(entry.Locations, loc)
		case extractedRegex.MatchString(t.String()):
			entry.ExtractedComments = append(entry.ExtractedComments,
				extractedRegex.FindStringSubmatch(t.String())[1],
			)
		case flagRegex.MatchString(t.String()):
			entry.Flags = append(entry.Flags,
				flagRegex.FindStringSubmatch(t.String())[1],
			)
		case previousRegex.MatchString(t.String()):
			entry.Previous = append(entry.Previous,
				previousRegex.FindStringSubmatch(t.String())[1],
			)
		default:
			entry.Comments = append(entry.Comments,
				generalRegex.FindStringSubmatch(t.String())[1],
			)
		}
	}

	return nil
}

// ParseWithOptions parses the PO file with temporary configuration options.
// The original configuration is restored after parsing.
func (p *PoParser) ParseWithOptions(opts ...PoOption) *po.File {
	p.Config.ApplyOptions(opts...)
	defer p.Config.RestoreLastCfg()

	return p.Parse()
}

const safeTemplate = `msgid "---"
msgstr "---"`

// Parse reads and parses the PO file into a po.File structure.
func (p *PoParser) Parse() *po.File {
	var entries po.Entries
	p.errors = nil

	// Add temporary marker to handle edge cases
	p.data = p.originalData
	p.data = append(p.data, []byte(safeTemplate)...)

	// Parse the file using the participle parser
	pFile, err := util.PoParser.ParseBytes(p.filename, p.data)
	if err != nil {
		p.error(err.Error())
		return nil
	}

	// Remove the temporary marker entry
	pFile.Entries = slices.Delete(pFile.Entries, len(pFile.Entries)-1, len(pFile.Entries))

	// Optionally filter out all comments
	if p.Config.IgnoreAllComments {
		pFile.Tokens = slices.DeleteFunc(pFile.Tokens, func(t lexer.Token) bool {
			return t.Type == util.PoSymbols["Comment"]
		})
	}

	// Process each entry in the parsed file
	for _, e := range pFile.Entries {
		if p.Config.IgnoreAllComments {
			e.Tokens = slices.DeleteFunc(e.Tokens, func(t lexer.Token) bool {
				return t.Type == util.PoSymbols["Comment"]
			})
		}

		// Create new entry with joined strings
		newEntry := po.Entry{
			Context:  strings.Join(e.Context, ""),
			ID:       strings.Join(e.ID, ""),
			Str:      strings.Join(e.Str, ""),
			Plural:   strings.Join(e.MsgidPlural, ""),
			Obsolete: p.Config.markAllAsObsolete,
		}

		// Process plural forms
		for _, pe := range e.Plurals {
			np := po.PluralEntry{
				ID:  pe.ID,
				Str: strings.Join(pe.Str, "\n"),
			}

			newEntry.Plurals = append(newEntry.Plurals, np)
		}

		// Parse comments associated with this entry
		err = parseComments(&newEntry, e.Tokens)
		if err != nil {
			p.error(err.Error())
		}

		// Optionally filter out comments
		if p.Config.IgnoreComments {
			newEntry.Comments = nil
			newEntry.ExtractedComments = nil
			newEntry.Previous = nil
		}

		entries = append(entries, newEntry)
	}

	// Process obsolete entries if configured to do so
	if slices.ContainsFunc(pFile.Tokens, func(t lexer.Token) bool {
		return t.Type == util.PoSymbols["Comment"] && obsoleteRegex.MatchString(t.String())
	}) && p.Config.ParseObsoletes {
		entries = append(entries, p.parseObsoleteEntries(pFile.Tokens)...)
	}

	// Apply post-processing options
	if p.Config.SkipHeader {
		i := entries.Index("", "")
		if i != -1 {
			entries = slices.Delete(entries, i, i+1)
		}
	}
	if p.Config.CleanDuplicates {
		entries = entries.CleanDuplicates()
	}

	return &po.File{
		Entries: entries,
		Name:    p.filename,
	}
}
