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
}

func New(input *ast.File, content []rune) *Generator {
	g := &Generator{file: input, content: content}
	return g
}

func (g *Generator) Generate() (f *types.File, warns []string, errs []error) {
	f = &types.File{
		Name: g.file.Name,
	}
	f.Entries, warns, errs = g.genEntries()
	f.Header = g.genHeader(f.LoadID(""))

	return
}

func (g *Generator) genHeader(str string) (header types.Header) {
	return
}

func tfor[T any]() reflect.Type {
	return reflect.TypeFor[T]()
}

func (g *Generator) genEntries() (entries []types.Entry, warns []string, errs []error) {
	var (
		curEntry types.Entry
		foundStr bool
		foundID  bool
	)

	reset := func() {
		curEntry = types.Entry{}
		foundStr = false
		foundID = false
	}

	finish := func(cur ast.Node) {
		err := validateEntry(curEntry)
		if err != nil {
			errs = append(errs, err)
			return
		}

		if !foundID {
			errs = append(errs,
				fmt.Errorf(
					"msgid not found at %s:%d",
					g.file.Name,
					util.FindLine(g.content, cur.Pos()),
				),
			)
			return
		}
		if !foundStr {
			warns = append(
				warns,
				fmt.Sprintf(
					"msgstr not found at %s:%d",
					g.file.Name,
					util.FindLine(g.content, cur.Pos()),
				),
			)
		}

		if !foundStr || !foundID {
			return
		}

		entries = append(entries, curEntry)
		reset()
	}

	toSkip := []reflect.Type{
		tfor[ast.FlagComment](),
		tfor[ast.LocationComment](),
		tfor[ast.GeneralComment](),
	}

	for i := 0; i < len(g.file.Nodes); i++ {
		node := g.file.Nodes[i]

		switch n := node.(type) {
		case ast.GeneralComment, ast.FlagComment, ast.LocationComment:
			switch n := n.(type) {
			case ast.LocationComment:
				loc := types.Location{
					File: n.File,
					Line: n.Line,
				}
				curEntry.Locations = append(curEntry.Locations, loc)
			case ast.FlagComment:
				curEntry.Flags = append(curEntry.Flags, n.Flag)
			default:
				continue
			}
		case ast.Msgctxt:
			if curEntry.Context != "" {
				warns = append(warns,
					fmt.Sprintf("duplicated msgctxt at %s:%d",
						g.file.Name,
						util.FindLine(g.content, n.Pos()),
					),
				)
			}
			curEntry.Context = n.Context

			if !g.typeIsComing(i+1,
				toSkip,
				tfor[ast.Msgid](),
			) {
				finish(n)
			}

		case ast.Msgid:
			if foundID {
				warns = append(warns,
					fmt.Sprintf("duplicated msgid at %s:%d",
						g.file.Name,
						util.FindLine(g.content, n.Pos()),
					),
				)
			}
			foundID = true
			curEntry.ID = n.ID

			if !g.typeIsComing(i+1,
				toSkip,
				tfor[ast.Msgstr](),
				tfor[ast.MsgidPlural](),
				tfor[ast.MsgstrPlural](),
			) {
				finish(n)
			}
		case ast.MsgidPlural:
			if curEntry.Plural != "" {
				warns = append(warns,
					fmt.Sprintf("duplicated msgid_plural at %s:%d",
						g.file.Name,
						util.FindLine(g.content, n.Pos()),
					),
				)
			}
			curEntry.Plural = n.Plural

			if !g.typeIsComing(i+1,
				toSkip,
				tfor[ast.MsgstrPlural](),
			) {
				finish(n)
			}
		case ast.MsgstrPlural:
			p := types.PluralEntry{
				ID:  n.PluralID,
				Str: n.Str,
			}
			curEntry.Plurals = append(curEntry.Plurals, p)
			foundStr = true
			if !g.typeIsComing(i+1, toSkip, tfor[ast.MsgstrPlural]()) {
				finish(n)
			}
		case ast.Msgstr:
			if curEntry.Plural != "" {
				continue
			}

			foundStr = true
			curEntry.Str = n.Str
			finish(n)
		}
	}
	return
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
