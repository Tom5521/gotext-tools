package compiler_test

import (
	"testing"

	"github.com/Tom5521/xgotext/pkg/po"
	"github.com/Tom5521/xgotext/pkg/po/compiler"
)

func BenchmarkMoCompiler(b *testing.B) {
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

	comp := compiler.NewMo(&po.File{Entries: input})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		comp.ToBytes()
	}
}
