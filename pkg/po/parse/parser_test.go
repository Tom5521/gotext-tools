// TODO: Remove repetitive code.
package parse_test

import (
	"fmt"
	"testing"

	"github.com/Tom5521/xgotext/pkg/po/parse"
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

	p := parse.NewParserFromString(input, "test.go")
	errs := p.Parse()
	if len(errs) > 0 {
		t.Errorf("Unexpected error: %v\n", errs[0])
		return
	}
	fmt.Println(parse.FormatNode(p.Nodes()...))

	expected := []parse.Node{
		parse.GeneralComment{
			Text: "General Comment",
		},
		parse.FlagComment{
			Flag: "flag comment",
		},
		parse.LocationComment{
			File: "location_comment",
			Line: 123,
		},
		parse.Msgctxt{
			Context: "testing!",
		},
		parse.Msgid{
			ID: "1st msgid!",
		},
		parse.Msgstr{
			Str: "1er msgid!",
		},
		parse.Msgid{
			ID: "I want an apple",
		},
		parse.MsgidPlural{
			Plural: "I want some apples",
		},
		parse.MsgstrPlural{
			PluralID: 0,
			Str:      "Quiero una manzana",
		},
		parse.MsgstrPlural{
			PluralID: 1,
			Str:      "Quiero unas manzanas",
		},
	}

	nodes := p.Nodes()

	if !parse.EqualNodeSlice(expected, nodes) {
		t.Error("Unexpected node slice...")
		t.Error("Expected:", expected)
		t.Error("Got:", nodes)
	}
}
