package parse

import (
	"io"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/Tom5521/xgotext/pkg/po"
	"github.com/alecthomas/participle/v2/lexer"
)

type PoParser struct {
	Config PoConfig

	data     []byte
	filename string

	errors []error
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
		Config:   DefaultPoConfig(options...),
		data:     data,
		filename: name,
	}
}

// Return the first error in the stack.
func (p PoParser) Error() error {
	if len(p.errors) == 0 {
		return nil
	}
	return p.errors[0]
}

func (p PoParser) Errors() []error {
	return p.errors
}

var (
	pluralIDRegex  = regexp.MustCompile(`\d+`)
	locationRegex  = regexp.MustCompile(`#: *(.*)`)
	generalRegex   = regexp.MustCompile(`# *(.*)`)
	extractedRegex = regexp.MustCompile(`#\. *(.*)`)
	flagRegex      = regexp.MustCompile(`#, *(.*)`)
	previousRegex  = regexp.MustCompile(`#\| *(.*)`)
)

func parseComments(entry *po.Entry, tks []lexer.Token) (err error) {
	for _, t := range tks {
		switch t.Type {
		case tokens["COMMENT"]:
			entry.Comments = append(entry.Comments,
				generalRegex.FindStringSubmatch(t.String())[1],
			)
		case tokens["FLAG_COMMENT"]:
			entry.Flags = append(entry.Flags,
				flagRegex.FindStringSubmatch(t.String())[1],
			)
		case tokens["EXTRACTED_COMMENT"]:
			entry.ExtractedComments = append(entry.ExtractedComments,
				extractedRegex.FindStringSubmatch(t.String())[1],
			)
		case tokens["PREVIOUS_COMMENT"]:
			entry.Previous = append(entry.Previous,
				previousRegex.FindStringSubmatch(t.String())[1],
			)
		case tokens["REFERENCE_COMMENT"]:
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
		default:
			continue
		}
	}

	return nil
}

func (p *PoParser) Parse(options ...PoOption) *po.File {
	p.Config.ApplyOptions(options...)
	defer p.Config.RestoreLastCfg()

	var entries po.Entries
	p.errors = nil

	pFile, err := poParser.ParseBytes(p.filename, p.data)
	if err != nil {
		p.Config.Logger.Println("ERROR:", err)
		p.errors = append(p.errors, err)
		return nil
	}

mainLoop:
	for _, e := range pFile.Entries {
		newEntry := po.Entry{
			Context: strings.Join(e.Context, "\n"),
			ID:      strings.Join(e.ID, "\n"),
			Str:     strings.Join(e.Str, "\n"),
			Plural:  strings.Join(e.MsgidPlural, "\n"),
		}

		// Parse plurals
		for _, pe := range e.Plurals {
			id, err := strconv.Atoi(pluralIDRegex.FindString(pe.ID))
			if err != nil {
				p.Config.Logger.Println("ERROR:", err)
				p.errors = append(p.errors, err)
				continue mainLoop
			}

			np := po.PluralEntry{
				ID:  id,
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

		entries = append(entries, newEntry)
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
