package compile_test

import (
	"fmt"
	"testing"

	"github.com/Tom5521/gotext-tools/v2/internal/util"
	"github.com/Tom5521/gotext-tools/v2/pkg/po"
	"github.com/Tom5521/gotext-tools/v2/pkg/po/compile"
	"github.com/Tom5521/gotext-tools/v2/pkg/po/parse"
)

func TestPoCompiler(t *testing.T) {
	input := po.Entries{
		{Context: "My context :3", ID: "id1", Str: "HELLO"},
		{
			ID:     "id2",
			Plural: "helooows\nLol",
			Plurals: po.PluralEntries{
				po.PluralEntry{ID: 0, Str: "Holanda"},
				po.PluralEntry{ID: 1, Str: "Holandas"},
			},
		},
		{ID: "id3", Str: "Hello3"},
	}

	tests := []struct {
		name    string
		options []compile.PoOption
	}{
		{"Base", nil},
		{"WithWordWrap", []compile.PoOption{compile.PoWithWordWrap(true)}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.options = append(test.options, compile.PoWithOmitHeader(true))
			compiled := compile.NewPo(&po.File{Entries: input}, test.options...).ToBytes()

			parser := parse.NewPoFromBytes(compiled, "test.po")

			parsed := parser.Parse().Entries
			if parser.Error() != nil {
				t.Error(parser.Error())
				fmt.Println(parser.Errors())
			}

			if !util.Equal(parsed, input) {
				t.Error("Input and output differ!")

				fmt.Println(util.NamedDiff("input", "output", input, parsed))
			}
		})
	}
}
