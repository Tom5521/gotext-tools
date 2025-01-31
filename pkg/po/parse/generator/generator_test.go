package generator_test

import (
	"reflect"
	"testing"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po/parse/ast"
	"github.com/Tom5521/xgotext/pkg/po/parse/generator"
	"github.com/Tom5521/xgotext/pkg/po/types"
	"github.com/kr/pretty"
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

	expected := types.Entries{
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

	g := generator.New(p.File)
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

func TestGenHeader(t *testing.T) {
	const input = `msgid ""
msgstr ""
"Project-Id-Version: PACKAGE VERSION\n"
"Report-Msgid-Bugs-To: \n"
"POT-Creation-Date: 2025-01-20 14:53:37\n"
"PO-Revision-Date: \n"
"Last-Translator: \n"
"Language-Team: \n"
"Language: en\n"
"MIME-Version: 1.0\n"
"Content-Type: text/plain; charset=CHARSET\n"
"Content-Transfer-Encoding: 8bit\n"
"Plural-Forms: nplurals=2; plural=(n != 1);\n"`

	p := ast.NewParserFromString(input, "test.go")
	errs := p.Parse()
	if errs != nil {
		t.Error("Unexpected error found:")
		t.Error(errs[0])
		return
	}

	g := generator.New(p.File)
	f := g.Generate()
	if len(g.Errors()) > 0 {
		t.Error("Unexpected error found:")
		t.Error(g.Errors()[0])
		return
	}
	for _, warn := range g.Warnings() {
		t.Log("WARN:", warn)
	}

	expected := &types.File{
		Name: "test.go",
		Header: types.Header{
			Values: []types.HeaderField{
				{Key: "Project-Id-Version", Value: "PACKAGE VERSION"},
				{Key: "Report-Msgid-Bugs-To", Value: ""},
				{Key: "POT-Creation-Date: 2025-01-20 14:53", Value: "37"},
				{Key: "PO-Revision-Date", Value: ""},
				{Key: "Last-Translator", Value: ""},
				{Key: "Language-Team", Value: ""},
				{Key: "Language", Value: "en"},
				{Key: "MIME-Version", Value: "1.0"},
				{Key: "Content-Type", Value: "text/plain; charset=CHARSET"},
				{Key: "Content-Transfer-Encoding", Value: "8bit"},
				{Key: "Plural-Forms", Value: "nplurals=2; plural=(n != 1);"},
			},
		},
		Nplurals: 2,
		Entries: types.Entries{
			{
				Flags:     []string(nil),
				ID:        "",
				Context:   "",
				Plural:    "",
				Plurals:   []types.PluralEntry(nil),
				Str:       "\nProject-Id-Version: PACKAGE VERSION\n\nReport-Msgid-Bugs-To: \n\nPOT-Creation-Date: 2025-01-20 14:53:37\n\nPO-Revision-Date: \n\nLast-Translator: \n\nLanguage-Team: \n\nLanguage: en\n\nMIME-Version: 1.0\n\nContent-Type: text/plain; charset=CHARSET\n\nContent-Transfer-Encoding: 8bit\n\nPlural-Forms: nplurals=2; plural=(n != 1);\n",
				Locations: []types.Location(nil),
			},
		},
	}

	if !reflect.DeepEqual(*expected, *f) {
		t.Error("Structures does not match:")
		t.Error("Expected:", pretty.Sprint(expected))
		t.Error("Got:", pretty.Sprint(f))
	}
}
