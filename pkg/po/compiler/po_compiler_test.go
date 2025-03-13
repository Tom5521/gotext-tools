package compiler_test

import (
	"fmt"
	"testing"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po"
	"github.com/Tom5521/xgotext/pkg/po/compiler"
	"github.com/Tom5521/xgotext/pkg/po/parse"
	"github.com/kr/pretty"
)

func TestPoCompiler(t *testing.T) {
	input := po.Entries{
		{Context: "My context :3", ID: "id1", Str: "HELLO"},
		{
			ID:     "id2",
			Plural: "helooows",
			Plurals: po.PluralEntries{
				po.PluralEntry{ID: 0, Str: "Holanda"},
				po.PluralEntry{ID: 1, Str: "Holandas"},
			},
		},
		{ID: "id3", Str: "Hello3"},
	}

	compiled := compiler.NewPo(&po.File{Entries: input}, compiler.PoWithOmitHeader(true)).ToBytes()

	parser := parse.NewPoFromBytes(compiled, "test.po")

	parsed := parser.Parse().Entries
	if parser.Error() != nil {
		t.Error(parser.Error())
		fmt.Println(parser.Errors())
	}

	if !util.Equal(parsed, input) {
		t.Error("Input and output differ!")
		fmt.Println("INPUT:\n", input)
		fmt.Println("OUTPUT:\n", parsed)
		fmt.Println("DIFF:")
		for _, d := range pretty.Diff(input, parsed) {
			fmt.Println(d)
		}
	}
}

func BenchmarkPoCompiler(b *testing.B) {
	input := po.Entries{
		{Context: "My context :3", ID: "id1", Str: "HELLO"},
		{
			ID:     "id2",
			Plural: "helooows",
			Plurals: po.PluralEntries{
				po.PluralEntry{ID: 0, Str: "Holanda"},
				po.PluralEntry{ID: 1, Str: "Holandas"},
			},
		},
		{ID: "id3", Str: "Hello3"},
	}

	comp := compiler.NewPo(&po.File{Entries: input})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		comp.ToBytes()
	}
}
