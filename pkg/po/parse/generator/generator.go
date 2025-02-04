package generator

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/Tom5521/xgotext/pkg/po/parse/ast"
	"github.com/Tom5521/xgotext/pkg/po/types"
)

type Generator struct {
	file *ast.File

	errs []error
}

func New(input *ast.File) *Generator {
	g := &Generator{
		file: input,
	}
	return g
}

func (g Generator) Errors() []error {
	return g.errs
}

func (g *Generator) Generate() (f *types.File) {
	f = &types.File{
		Name:    g.file.Name,
		Entries: g.genEntries(),
	}
	g.genHeader(f)
	g.genNplurals(f)

	return
}

var (
	npluralsRegex = regexp.MustCompile(`nplurals=(\d*)`)
	headerRegex   = regexp.MustCompile(`(.*)\s*:\s*(.*)`)
)

func (g *Generator) genHeader(f *types.File) {
	header := f.LoadID("")
	lines := strings.Split(header, "\n")
	for _, line := range lines {
		if !headerRegex.MatchString(line) {
			continue
		}
		matches := headerRegex.FindStringSubmatch(line)
		f.Header.Values = append(f.Header.Values,
			types.HeaderField{
				Key:   matches[1],
				Value: matches[2],
			},
		)
	}
}

func (g *Generator) genNplurals(f *types.File) {
	for _, field := range f.Header.Values {
		if field.Key == "Plural-Forms" {
			if !npluralsRegex.MatchString(field.Value) {
				break
			}
			matches := npluralsRegex.FindStringSubmatch(field.Value)
			nplurals, err := strconv.Atoi(matches[1])
			if err != nil {
				g.errs = append(g.errs, err)
				break
			}
			f.Nplurals = nplurals
			break
		}
	}
}

func (g *Generator) genEntries() types.Entries {
	var entries types.Entries
	for _, node := range g.file.Nodes {
		raw, ok := node.(ast.Entry)
		if !ok {
			continue
		}

		var entry types.Entry

		if raw.Msgid != nil {
			entry.ID = raw.Msgid.ID
		}
		if raw.Msgstr != nil {
			entry.Str = raw.Msgstr.Str
		}
		if raw.Msgctxt != nil {
			entry.Context = raw.Msgctxt.Context
		}
		if raw.Plural != nil {
			entry.Plural = raw.Plural.Plural
		}

		for _, flag := range raw.Flags {
			entry.Flags = append(entry.Flags, flag.Flag)
		}
		for _, extractedComment := range raw.ExtractedComments {
			entry.ExtractedComments = append(entry.ExtractedComments, extractedComment.Text)
		}
		for _, location := range raw.LocationComments {
			entry.Locations = append(
				entry.Locations,
				types.Location{Line: location.Line, File: location.File},
			)
		}
		for _, generalComment := range raw.GeneralComments {
			entry.Comments = append(entry.Comments, generalComment.Text)
		}
		for _, plural := range raw.Plurals {
			entry.Plurals = append(entry.Plurals, types.PluralEntry{
				ID:  plural.PluralID,
				Str: plural.Str,
			})
		}

		entries = append(entries, entry)
	}

	return entries
}
