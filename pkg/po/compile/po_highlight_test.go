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
	highlighted, err := compile.HighlighOutput(
		compile.DefaultHighligh,
		"test.po",
		bytes.NewBufferString(input),
	)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(highlighted)
}
