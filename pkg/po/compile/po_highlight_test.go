package compile_test

import (
	"testing"

	"github.com/Tom5521/gotext-tools/v2/pkg/po/compile"
)

func TestHighlight(t *testing.T) {
	input := `# This is a comment
msgctxt "WAOS"
msgid "Lol"
msgstr "waos"
msgstr[1] "LOL"`

	highlighted, err := compile.Highlight(compile.DefaultHighlight, "input.po", input)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(highlighted))
}

func BenchmarkHighlight(b *testing.B) {
	input := `# This is a comment
msgctxt "WAOS"
msgid "Lol"
msgstr "waos"
msgstr[1] "LOL"`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := compile.Highlight(compile.DefaultHighlight, "input.po", input)
		b.StopTimer()
		if err != nil {
			b.Error(err)
		}
		b.StartTimer()
	}
}
