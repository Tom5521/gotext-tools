// TODO: Remove repetitive code.
package ast_test

import (
	"testing"

	"github.com/Tom5521/xgotext/pkg/po/parse/ast"
)

func TestParse(t *testing.T) {
	const input = `# General Comment
#, flag comment
#: location_comment:123
msgctxt "testing!"
msgid "1st msgid!"
msgstr "1er msgid!"

msgid "I want an apple"
msgid_plural "I want some apples"
msgstr[0] "Quiero una manzana"
msgstr[1] "Quiero unas manzanas"`

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
	}

	nodes := p.Nodes()

	if !ast.EqualNodeSlice(expected, nodes) {
		t.Error("Unexpected node slice...")
		t.Error("Expected:", ast.FormatNode(expected...))
		t.Error("Got:", ast.FormatNode(nodes...))
	}
}
