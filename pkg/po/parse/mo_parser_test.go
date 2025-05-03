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

func TestMoParse(t *testing.T) {
	entries := po.Entries{
		{
			ID:      "Apple",
			Context: "USA",
			Plural:  "Apples",
			Plurals: po.PluralEntries{
				{ID: 0, Str: "Manzana"},
				{ID: 1, Str: "Manzanas"},
			},
		},
		{ID: "Hi", Str: "Hola", Context: "casual"},
		{ID: "How are you?", Str: "Como estás?"},
	}.SortFunc(po.CompareEntryByID)

	moFile := compile.MoToBytes(entries, compile.MoWithHashTable(false))
	parser := parse.NewMoFromBytes(moFile, "test.mo")

	parsedEntries := parser.Parse()
	if parser.Error() != nil {
		t.Error(parser.Error())
		return
	}

	if !util.Equal(entries, parsedEntries.Entries) {
		t.Error("Parsed entries differ!")
		fmt.Println("--- ORIGINAL:", entries)
		fmt.Println("--- PARSED:", parsedEntries.Entries)
		fmt.Println("--- DIFF:")
		d := diff.Diff(
			"good",
			compile.PoToBytes(entries, compile.PoWithOmitHeader(true)),
			"bad",
			compile.PoToBytes(parsedEntries, compile.PoWithOmitHeader(true)),
		)

		fmt.Println(string(d))
		return
	}
}
