package parse_test

import (
	"fmt"
	"testing"

	"github.com/Tom5521/gotext-tools/v2/internal/util"
	"github.com/Tom5521/gotext-tools/v2/pkg/po"
	"github.com/Tom5521/gotext-tools/v2/pkg/po/compile"
	"github.com/Tom5521/gotext-tools/v2/pkg/po/parse"
	"github.com/rogpeppe/go-internal/diff"
)

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

	comp := compile.NewPo(input, compile.PoWithOmitHeader(true))
	expected := comp.ToString()

	parser := parse.NewPoFromString(expected, "test.po")
	parsed := parser.Parse()

	if parser.Error() != nil {
		t.Error(parser.Error())
		return
	}

	if !util.Equal(parsed.Entries, input.Entries) {
		t.Error("Compiled and parsed differ!")

		d := diff.Diff("parsed", compile.PoToBytes(parsed), "expected", compile.PoToBytes(input))
		fmt.Println(string(d))
	}
}
