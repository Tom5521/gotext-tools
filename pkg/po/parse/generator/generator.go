package generator

import (
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
	return types.NewFile(g.file.Name, g.genEntries()...)
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
