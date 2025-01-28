package generator_test

import (
	"testing"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po/parse/ast"
	"github.com/Tom5521/xgotext/pkg/po/parse/generator"
	"github.com/Tom5521/xgotext/pkg/po/types"
)

func TestGen(t *testing.T) {
	const input = `# hello.go:123
msgid "Hi"
msgstr "Hola"

#, myflag
#: myfile:12
msgctxt "formal"
msgid "Hello"
msgstr "Saludos"

#, flag1
#: Hello.go:123
msgid "You have %d apple"
msgid_plural "You have %d apples"
msgstr[0] "Tienes %d manzana"
msgstr[1] "Tienes %d manzanas"`

	expected := []types.Entry{
		{
			Flags:     nil,
			ID:        "Hi",
			Context:   "",
			Plural:    "",
			Plurals:   nil,
			Str:       "Hola",
			Locations: nil,
		}, {
			Flags:   []string{"myflag"},
			ID:      "Hello",
			Context: "formal",
			Plural:  "",
			Plurals: nil,
			Str:     "Saludos",
			Locations: []types.Location{
				{Line: 12, File: "myfile"},
			},
		}, {
			Flags:   []string{"flag1"},
			ID:      "You have %d apple",
			Context: "",
			Plural:  "You have %d apples",
			Plurals: []types.PluralEntry{
				{ID: 0, Str: "Tienes %d manzana"},
				{ID: 1, Str: "Tienes %d manzanas"},
			},
			Str: "",
			Locations: []types.Location{
				{Line: 123, File: "Hello.go"},
			},
		},
	}

	p := ast.NewParserFromString(input, "test.go")
	p.Parse()

	g := generator.New(p.File, []rune(input))
	file := g.Generate()
	if len(g.Errors()) > 0 {
		t.Error("Unexpected error found:")
		t.Error(g.Errors()[0])
		return
	}
	for _, warn := range g.Warnings() {
		t.Log("WARN:", warn)
	}

	if !types.EqualEntries(expected, file.Entries) {
		t.Error("The results does not match")
		t.Error("Expected", util.Format(expected))
		t.Error("Got:", util.Format(file.Entries))
	}
}
