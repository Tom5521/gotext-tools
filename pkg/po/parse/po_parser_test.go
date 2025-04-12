package parse_test

import (
	"fmt"
	"testing"

	"github.com/Tom5521/gotext-tools/internal/util"
	"github.com/Tom5521/gotext-tools/pkg/po"
	"github.com/Tom5521/gotext-tools/pkg/po/compiler"
	"github.com/Tom5521/gotext-tools/pkg/po/parse"
	"github.com/sergi/go-diff/diffmatchpatch"
)

var dmp = diffmatchpatch.New()

func TestPoParser(t *testing.T) {
	input := &po.File{
		Entries: po.Entries{
			{
				Flags:    []string{"my-flag lol"},
				Comments: []string{"Hello World"},
				ID:       "Hello", Str: "Hola",
			},
			{Context: "CTX", ID: "MEOW", Str: "MIAU"},
			{
				ID:      "Apple",
				Plural:  "Apples",
				Plurals: po.PluralEntries{{ID: 0, Str: "Manzana"}, {ID: 1, Str: "Manzanas"}},
			},
			{
				ID:       "MyObsoleteEntry",
				Obsolete: true,
			},
		},
	}

	comp := compiler.NewPo(input, compiler.PoWithOmitHeader(true))
	expected := comp.ToString()

	parser := parse.NewPoFromString(expected, "test.po")
	parsed := parser.Parse()

	if parser.Error() != nil {
		t.Error(parser.Error())
		return
	}

	if !util.Equal(parsed.Entries, input.Entries) {
		t.Error("Compiled and parsed differ!")

		comp.File = parsed

		dmain := dmp.DiffMain(comp.ToString(), expected, false)
		fmt.Println(dmp.DiffPrettyText(dmain))
	}
}
