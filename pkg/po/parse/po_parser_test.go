package parse_test

import (
	"fmt"
	"testing"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po"
	"github.com/Tom5521/xgotext/pkg/po/compiler"
	"github.com/Tom5521/xgotext/pkg/po/parse"
	"github.com/kr/pretty"
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
		},
	}

	compiled := compiler.NewPo(input, compiler.PoWithOmitHeader(true)).ToString()

	parser := parse.NewPoFromString(compiled, "test.po")
	parsed := parser.Parse()

	if parser.Error() != nil {
		t.Error(parser.Error())
		return
	}

	if !util.Equal(parsed.Entries, input.Entries) {
		t.Error("Compiled and parsed differ!")
		fmt.Println("ORIGINAL:\n", compiled)
		fmt.Println("PARSED:\n", compiler.NewPo(parsed, compiler.PoWithOmitHeader(true)).ToString())

		fmt.Println("DIFF:")
		for _, d := range pretty.Diff(parsed.Entries, input.Entries) {
			fmt.Println(d)
		}
	}
}

func BenchmarkParsePo(b *testing.B) {
	input := `# 
# Copyright (C) 
# This file is distributed under the same license as the PACKAGE NAME package.
#
msgid ""
msgstr ""

msgctxt "My context :3"
msgid "id1"
msgstr "id1"

msgid "id2"
msgid_plural "helooows"
msgstr[0] "Holanda"
msgstr[1] "Holandas"

msgid "id3"
msgstr "id3"`

	parser := parse.NewPoFromString(input, "test.po")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser.Parse()
		b.StopTimer()
		if err := parser.Error(); err != nil {
			b.Error(err)
		}
		b.StartTimer()
	}
}
