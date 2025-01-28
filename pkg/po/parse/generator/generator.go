package generator

import (
	"errors"
	"fmt"
	"reflect"
	"slices"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po/parse/ast"
	"github.com/Tom5521/xgotext/pkg/po/types"
)

type Generator struct {
	content []rune
	file    *ast.File

	header types.Header

	curEntry types.Entry
	foundStr bool
	foundID  bool
	entries  []types.Entry
	warns    []string
	errs     []error

	toSkip []reflect.Type
}

func New(input *ast.File, content []rune) *Generator {
	g := &Generator{
		file:    input,
		content: content,
		toSkip: []reflect.Type{
			tfor[ast.FlagComment](),
			tfor[ast.GeneralComment](),
			tfor[ast.LocationComment](),
		},
	}
	return g
}

func (g Generator) Errors() []error {
	return g.errs
}

func (g Generator) Warnings() []string {
	return g.warns
}

func (g *Generator) Generate() (f *types.File) {
	f = &types.File{
		Name: g.file.Name,
	}
	g.genEntries()
	g.genHeader(f.LoadID(""))

	f.Entries = g.entries
	f.Header = g.header

	return
}

func (g *Generator) genHeader(str string) {
}

func tfor[T any]() reflect.Type {
	return reflect.TypeFor[T]()
}

func (g *Generator) genEntries() {
	g.resetEntryState()

	for i, node := range g.file.Nodes {
		switch n := node.(type) {
		case ast.GeneralComment, ast.FlagComment, ast.LocationComment:
			g.handleComment(n)
		case ast.Msgctxt:
			g.handleMsgctxt(n, i)
		case ast.Msgid:
			g.handleMsgid(n, i)
		case ast.MsgidPlural:
			g.handleMsgidPlural(n, i)
		case ast.MsgstrPlural:
			g.handleMsgstrPlural(n, i)
		case ast.Msgstr:
			g.handleMsgstr(n)
		}
	}
}

func (g *Generator) resetEntryState() {
	g.curEntry = types.Entry{}
	g.foundStr = false
	g.foundID = false
}

func (g *Generator) handleComment(node ast.Node) {
	switch n := node.(type) {
	case ast.LocationComment:
		loc := types.Location{
			File: n.File,
			Line: n.Line,
		}
		g.curEntry.Locations = append(g.curEntry.Locations, loc)
	case ast.FlagComment:
		g.curEntry.Flags = append(g.curEntry.Flags, n.Flag)
	default:
		return
	}
}

func (g *Generator) handleMsgctxt(n ast.Msgctxt, i int) {
	if g.curEntry.Context != "" {
		g.warns = append(g.warns,
			fmt.Sprintf("duplicated msgctxt at %s:%d",
				g.file.Name,
				util.FindLine(g.content, n.Pos()),
			),
		)
	}
	g.curEntry.Context = n.Context

	if !g.typeIsComing(i+1,
		g.toSkip,
		tfor[ast.Msgid](),
	) {
		g.finishEntry(n)
	}
}

func (g *Generator) handleMsgid(n ast.Msgid, i int) {
	if g.foundID {
		g.warns = append(g.warns,
			fmt.Sprintf("duplicated msgid at %s:%d",
				g.file.Name,
				util.FindLine(g.content, n.Pos()),
			),
		)
	}
	g.foundID = true
	g.curEntry.ID = n.ID

	if !g.typeIsComing(i+1,
		g.toSkip,
		tfor[ast.Msgstr](),
		tfor[ast.MsgidPlural](),
		tfor[ast.MsgstrPlural](),
	) {
		g.finishEntry(n)
	}
}

func (g *Generator) handleMsgidPlural(n ast.MsgidPlural, i int) {
	if g.curEntry.Plural != "" {
		g.warns = append(g.warns,
			fmt.Sprintf("duplicated msgid_plural at %s:%d",
				g.file.Name,
				util.FindLine(g.content, n.Pos()),
			),
		)
	}
	g.curEntry.Plural = n.Plural

	if !g.typeIsComing(i+1,
		g.toSkip,
		tfor[ast.MsgstrPlural](),
	) {
		g.finishEntry(n)
	}
}

func (g *Generator) handleMsgstrPlural(n ast.MsgstrPlural, i int) {
	p := types.PluralEntry{
		ID:  n.PluralID,
		Str: n.Str,
	}
	g.curEntry.Plurals = append(g.curEntry.Plurals, p)
	g.foundStr = true
	if !g.typeIsComing(i+1, g.toSkip, tfor[ast.MsgstrPlural]()) {
		g.finishEntry(n)
	}
}

func (g *Generator) handleMsgstr(n ast.Msgstr) {
	if g.curEntry.Plural != "" {
		return
	}

	g.foundStr = true
	g.curEntry.Str = n.Str
	g.finishEntry(n)
}

func (g *Generator) finishEntry(cur ast.Node) {
	err := validateEntry(g.curEntry)
	if err != nil {
		g.errs = append(g.errs, err)
		return
	}

	if !g.foundID {
		g.errs = append(g.errs,
			fmt.Errorf(
				"msgid not found at %s:%d",
				g.file.Name,
				util.FindLine(g.content, cur.Pos()),
			),
		)
		return
	}
	if !g.foundStr {
		g.warns = append(
			g.warns,
			fmt.Sprintf(
				"msgstr not found at %s:%d",
				g.file.Name,
				util.FindLine(g.content, cur.Pos()),
			),
		)
	}

	if !g.foundStr || !g.foundID {
		return
	}

	g.entries = append(g.entries, g.curEntry)
	g.resetEntryState()
}

func (g *Generator) commingType(offset int, ignore []reflect.Type) reflect.Type {
	for _, node := range g.file.Nodes[offset:] {
		t := reflect.TypeOf(node)
		if slices.Contains(ignore, t) {
			continue
		}

		return t
	}

	return nil
}

func (g *Generator) typeIsComing(offset int, ignore []reflect.Type, wanted ...reflect.Type) bool {
	t := g.commingType(offset, ignore)
	return slices.Contains(wanted, t)
}

func validateEntry(entry types.Entry) error {
	if entry.Plural == "" && len(entry.Plural) > 0 {
		return errors.New("plural translations provided but no plural form has been specified")
	}

	if entry.Plural != "" && len(entry.Plurals) == 0 {
		return errors.New("plural form specified but no plural translations provided")
	}

	return nil
}
