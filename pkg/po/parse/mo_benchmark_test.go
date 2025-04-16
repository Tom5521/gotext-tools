package parse_test

import (
	"bytes"
	"testing"

	"github.com/Tom5521/gotext-tools/v2/pkg/po"
	"github.com/Tom5521/gotext-tools/v2/pkg/po/compile"
	"github.com/Tom5521/gotext-tools/v2/pkg/po/parse"
)

func BenchmarkParseMo(b *testing.B) {
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
		{ID: "", Str: ""},
		{ID: "How are you?", Str: "Como estás?"},
	}

	compiled := compile.NewMo(&po.File{Entries: entries}).ToBytes()

	reader := bytes.NewReader(compiled)
	parser, err := parse.NewMoFromReader(reader, "test.mo")
	if err != nil {
		b.Error(err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser.Parse()
		b.StopTimer()
		if parser.Error() != nil {
			b.Error(parser.Error())
			b.Skip(parser.Error())
		}
		reader.Reset(compiled)
		b.StartTimer()
	}
}
