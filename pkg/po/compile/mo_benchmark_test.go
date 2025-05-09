package compile_test

import (
	"testing"

	"github.com/Tom5521/gotext-tools/v2/pkg/po"
	"github.com/Tom5521/gotext-tools/v2/pkg/po/compile"
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

	benchmarks := []struct {
		name string
		opts []compile.MoOption
	}{
		{
			"WithHashTable",
			[]compile.MoOption{compile.MoWithHashTable(true)},
		},
		{
			"WithoutHashTable",
			[]compile.MoOption{compile.MoWithHashTable(false)},
		},
	}

	for _, bench := range benchmarks {
		b.Run(bench.name, func(b *testing.B) {
			comp := compile.NewMo(&po.File{Entries: input}, bench.opts...)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				comp.ToBytes()
			}
		})
	}
}
