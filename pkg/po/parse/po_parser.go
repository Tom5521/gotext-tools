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

var _ po.Parser = (*PoParser)(nil)

type PoParser struct {
	Config PoConfig

	originalData []byte
	data         []byte
	filename     string

	errors []error
	warns  []error
}

// This function MUST be used to log any errors inside this structure.
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

func NewPo(path string, options ...PoOption) (*PoParser, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return NewPoFromBytes(file, path, options...), nil
}

func NewPoFromReader(r io.Reader, name string, options ...PoOption) (*PoParser, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return NewPoFromBytes(data, name, options...), nil
}

func NewPoFromFile(f *os.File, options ...PoOption) (*PoParser, error) {
	return NewPoFromReader(f, f.Name(), options...)
}

func NewPoFromString(s, name string, options ...PoOption) *PoParser {
	return NewPoFromBytes([]byte(s), name, options...)
}

func NewPoFromBytes(data []byte, name string, options ...PoOption) *PoParser {
	return &PoParser{
		Config:       DefaultPoConfig(options...),
		originalData: data,
		filename:     name,
	}
}

// Return the first error in the stack.
func (p PoParser) Error() error {
	if len(p.errors) == 0 {
		return nil
	}
	return p.errors[0]
}

func (p PoParser) Warnings() []error {
	return p.warns
}

func (p PoParser) Errors() []error {
	return p.errors
}

var (
	locationRegex  = regexp.MustCompile(`#: *(.*)`)
	generalRegex   = regexp.MustCompile(`# *(.*)`)
	extractedRegex = regexp.MustCompile(`#\. *(.*)`)
	flagRegex      = regexp.MustCompile(`#, *(.*)`)
	obsoleteRegex  = regexp.MustCompile(`#~ *(.*)`)
	previousRegex  = regexp.MustCompile(`#\| *(.*)`)
)

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

	var cleanedLines []string
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
					return err
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

func (p *PoParser) ParseWithOptions(opts ...PoOption) *po.File {
	p.Config.ApplyOptions(opts...)
	defer p.Config.RestoreLastCfg()

	return p.Parse()
}

func (p *PoParser) Parse() *po.File {
	var entries po.Entries
	p.errors = nil

	p.data = p.originalData
	// Obsolete securer.
	p.data = append(p.data, []byte(`msgid "---"
msgstr "---"`)...)

	pFile, err := util.PoParser.ParseBytes(p.filename, p.data)
	if err != nil {
		p.error(err.Error())
		return nil
	}

	// Remove obsolete securer.
	pFile.Entries = slices.Delete(pFile.Entries, len(pFile.Entries)-1, len(pFile.Entries))

	if p.Config.IgnoreAllComments {
		pFile.Tokens = slices.DeleteFunc(pFile.Tokens, func(t lexer.Token) bool {
			return t.Type == util.PoSymbols["Comment"]
		})
	}

	for _, e := range pFile.Entries {
		if p.Config.IgnoreAllComments {
			e.Tokens = slices.DeleteFunc(e.Tokens, func(t lexer.Token) bool {
				return t.Type == util.PoSymbols["Comment"]
			})
		}

		newEntry := po.Entry{
			Context:  strings.Join(e.Context, "\n"),
			ID:       strings.Join(e.ID, "\n"),
			Str:      strings.Join(e.Str, "\n"),
			Plural:   strings.Join(e.MsgidPlural, "\n"),
			Obsolete: p.Config.markAllAsObsolete,
		}

		// Parse plurals
		for _, pe := range e.Plurals {
			np := po.PluralEntry{
				ID:  pe.ID,
				Str: strings.Join(pe.Str, "\n"),
			}

			newEntry.Plurals = append(newEntry.Plurals, np)
		}
		// Parse Comments.
		err = parseComments(&newEntry, e.Tokens)
		if err != nil {
			p.error(err.Error())
		}

		if p.Config.IgnoreComments {
			newEntry.Comments = nil
			newEntry.ExtractedComments = nil
			newEntry.Previous = nil
		}

		entries = append(entries, newEntry)
	}

	if slices.ContainsFunc(pFile.Tokens, func(t lexer.Token) bool {
		return t.Type == util.PoSymbols["Comment"] && obsoleteRegex.MatchString(t.String())
	}) && p.Config.ParseObsoletes {
		entries = append(entries, p.parseObsoleteEntries(pFile.Tokens)...)
	}

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
