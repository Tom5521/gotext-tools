// TODO: Remove repetitive code.
package ast_test

import (
	"testing"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po/parse/ast"
	"github.com/kr/pretty"
)

func TestParse(t *testing.T) {
	const input = `# General Comment
#, flag comment
#: location_comment:123
#. extracted comment
msgctxt "testing!"
msgid "1st msgid!"
msgstr "1er msgid!"

msgid "I want an apple"
msgid_plural "I want some apples"
msgstr[0] "Quiero una manzana"
msgstr[1] "Quiero unas manzanas"

msgid ""
"hello"
"world"
""
msgstr ""
"hola"
"mundo"
""`

	p := ast.NewParserFromString(input, "test.go")
	errs := p.Parse()
	if len(errs) > 0 {
		t.Errorf("Unexpected error: %v\n", errs[0])
		return
	}

	expected := []ast.Node{
		ast.GeneralComment{
			Text: "General Comment",
		},
		ast.FlagComment{
			Flag: "flag comment",
		},
		ast.LocationComment{
			File: "location_comment",
			Line: 123,
		},
		ast.ExtractedComment{
			Text: "extracted comment",
		},
		ast.Msgctxt{
			Context: "testing!",
		},
		ast.Msgid{
			ID: "1st msgid!",
		},
		ast.Msgstr{
			Str: "1er msgid!",
		},
		ast.Msgid{
			ID: "I want an apple",
		},
		ast.MsgidPlural{
			Plural: "I want some apples",
		},
		ast.MsgstrPlural{
			PluralID: 0,
			Str:      "Quiero una manzana",
		},
		ast.MsgstrPlural{
			PluralID: 1,
			Str:      "Quiero unas manzanas",
		},
		ast.Msgid{
			ID: "\nhello\nworld\n",
		},
		ast.Msgstr{
			Str: "\nhola\nmundo\n",
		},
	}

	nodes := p.Nodes()

	if !util.Equal(expected, nodes) {
		t.Error("Unexpected node slice...")
		t.Error("Expected:", pretty.Sprint(expected))
		t.Error("Got:", pretty.Sprint(nodes))

		for _, d := range pretty.Diff(nodes, expected) {
			t.Error(d)
		}
	}
}
