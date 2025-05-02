package compile_test

import (
	"bytes"
	"testing"

	"github.com/Tom5521/gotext-tools/v2/pkg/po/compile"
)

func TestHighlight(t *testing.T) {
	input := `# This is a comment
msgctxt "WAOS"
msgid "Lol"
msgstr "waos"
msgstr[1] "LOL"`
	highlighted, err := compile.HighlightOutput(
		compile.DefaultHighlight,
		"test.po",
		bytes.NewBufferString(input),
	)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(highlighted))
}

func BenchmarkHighlight(b *testing.B) {
	str := `# This is a comment
msgctxt "WAOS"
msgid "Lol"
msgstr "waos"
msgstr[1] "LOL"`
	input := bytes.NewReader([]byte(str))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		compile.HighlightOutput(compile.DefaultHighlight, "input.po", input)
	}
}
