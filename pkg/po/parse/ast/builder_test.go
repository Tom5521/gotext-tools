package ast_test

import (
	"testing"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po/parse/ast"
	"github.com/kr/pretty"
)

func TestASTBuilder(t *testing.T) {
	input := `# general comment
#. extracted comment
#, flag comment
#: reference:123
msgid "Hi"
msgstr "Hola"

msgctxt "formal"
msgid "Hi"
msgstr "Buenas"

msgid "I want %d apple"
msgid_plural "I want %d apples"
msgstr[0] "Quiero %d manzana"
msgstr[1] "Quiero %d manzanas"`

	tokenizer := ast.NewTokenizer([]byte(input), "test.po")
	normalizer, errs := tokenizer.Normalizer()
	if len(errs) > 0 {
		t.Error(errs[0])
		return
	}

	normalizer.Build()

	if len(normalizer.Errors()) > 0 {
		t.Error(normalizer.Errors()[0])
		return
	}

	for _, warn := range normalizer.Warnings() {
		t.Log(warn)
	}

	expected := []ast.Entry{
		{
			Flags: []*ast.FlagComment{
				{Flag: "flag comment"},
			},
			ExtractedComments: []*ast.ExtractedComment{
				{Text: "extracted comment"},
			},
			LocationComments: []*ast.LocationComment{
				{File: "reference", Line: 123},
			},
			GeneralComments: []*ast.GeneralComment{
				{Text: "general comment"},
			},
			Msgid:   &ast.Msgid{ID: "Hi"},
			Msgstr:  &ast.Msgstr{Str: "Hola"},
			Msgctxt: nil,
			Plural:  nil,
			Plurals: nil,
		},
		{
			Flags:             nil,
			ExtractedComments: nil,
			LocationComments:  nil,
			GeneralComments:   nil,
			Msgid:             &ast.Msgid{ID: "Hi"},
			Msgstr:            &ast.Msgstr{Str: "Buenas"},
			Msgctxt:           &ast.Msgctxt{Context: "formal"},
			Plural:            nil,
			Plurals:           nil,
		},
		{
			Flags:             nil,
			ExtractedComments: nil,
			LocationComments:  nil,
			GeneralComments:   nil,
			Msgid:             &ast.Msgid{ID: "I want %d apple"},
			Msgstr:            nil,
			Msgctxt:           nil,
			Plural:            &ast.MsgidPlural{Plural: "I want %d apples"},
			Plurals: []*ast.MsgstrPlural{
				{PluralID: 0, Str: "Quiero %d manzana"},
				{PluralID: 1, Str: "Quiero %d manzanas"},
			},
		},
	}

	if !util.Equal(expected, normalizer.Entries()) {
		for _, d := range pretty.Diff(expected, normalizer.Entries()) {
			t.Error(d)
		}
	}
}
