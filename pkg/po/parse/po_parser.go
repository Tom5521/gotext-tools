package parse

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/Tom5521/gotext-tools/pkg/po"
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

// Parse directly the provided file.
func ParsePo(path string, opts ...PoOption) (*po.File, error) {
	parser, err := NewPo(path, opts...)
	if err != nil {
		return nil, err
	}
	file := parser.Parse()

	return file, parser.Error()
}

func ParsePoFromReader(r io.Reader, name string, opts ...PoOption) (*po.File, error) {
	parser, err := NewPoFromReader(r, name, opts...)
	if err != nil {
		return nil, err
	}

	file := parser.Parse()
	return file, parser.Error()
}

func ParsePoFromFile(f *os.File, opts ...PoOption) (*po.File, error) {
	parser, err := NewPoFromFile(f, opts...)
	if err != nil {
		return nil, err
	}
	file := parser.Parse()
	return file, parser.Error()
}

func ParsePoFromString(s, name string, opts ...PoOption) (*po.File, error) {
	parser := NewPoFromString(s, name, opts...)

	file := parser.Parse()
	return file, parser.Error()
}

func ParsePoFromBytes(b []byte, name string, opts ...PoOption) (*po.File, error) {
	parser := NewPoFromBytes(b, name, opts...)
	file := parser.Parse()

	return file, parser.Error()
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
		comp, err = regexp.Compile(fmt.Sprintf(`#%v *(.*)`, p.Config.CustomObsoletePrefix))
		if err != nil {
			p.warns = append(p.warns, fmt.Errorf("WARN: %w", err))
			return nil
		}

	}

	var cleanedLines []string
	for _, token := range tokens {
		if token.Type != symbols["Comment"] {
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
		for _, err2 := range parser.Errors() {
			err2 = fmt.Errorf("WARN: %w", err2)
			p.warns = append(p.warns, err2)
		}
		return nil
	}

	return f.Entries
}

func parseComments(entry *po.Entry, tokens []lexer.Token) (err error) {
	for _, t := range tokens {
		if t.Type != symbols["Comment"] {
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

	pFile, err := poParser.ParseBytes(p.filename, p.data)
	if err != nil {
		p.Config.Logger.Println("ERROR:", err)
		p.errors = append(p.errors, err)
		return nil
	}

	// Remove obsolete securer.
	pFile.Entries = slices.Delete(pFile.Entries, len(pFile.Entries)-1, len(pFile.Entries))

	if p.Config.IgnoreAllComments {
		pFile.Tokens = slices.DeleteFunc(pFile.Tokens, func(t lexer.Token) bool {
			return t.Type == symbols["Comment"]
		})
	}

	for _, e := range pFile.Entries {
		if p.Config.IgnoreAllComments {
			e.Tokens = slices.DeleteFunc(e.Tokens, func(t lexer.Token) bool {
				return t.Type == symbols["Comment"]
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
			p.Config.Logger.Println("ERROR:", err)
			p.errors = append(p.errors, err)
		}

		if p.Config.IgnoreComments {
			newEntry.Comments = nil
			newEntry.ExtractedComments = nil
			newEntry.Previous = nil
		}

		entries = append(entries, newEntry)
	}

	if slices.ContainsFunc(pFile.Tokens, func(t lexer.Token) bool {
		return t.Type == symbols["Comment"] && obsoleteRegex.MatchString(t.String())
	}) && p.Config.ParseObsoletes {
		entries = append(entries, p.parseObsoleteEntries(pFile.Tokens)...)
	}

	for _, err := range p.errors {
		p.Config.Logger.Println("ERROR:", err)
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
