package parse_test

import (
	"testing"

	. "github.com/Tom5521/xgotext/pkg/po/parse"
)

func TestParse(t *testing.T) {
	input := `#: file:32
msgid "MEOW!"
msgstr "LOL"
msgctxt "WOAS"
msgid "MEOW!"
msgstr "MIAU!"`

	l := NewLexer([]rune(input))

	expectedTokens := []Token{
		{COMMENT, ": file:32"},
		{MSGID, "msgid"},
		{STRING, "MEOW!"},
		{MSGSTR, "msgstr"},
		{STRING, "LOL"},
		{MSGCTXT, "msgctxt"},
		{STRING, "WOAS"},
		{MSGID, "msgid"},
		{STRING, "MEOW!"},
		{MSGSTR, "msgstr"},
		{STRING, "MIAU!"},
	}

	for i, etok := range expectedTokens {
		ctok := l.NextToken()

		if etok.Literal != ctok.Literal && etok.Type != ctok.Type {
			t.Errorf("unexpected token [%d]:", i)
			t.Error("got:", ctok)
			t.Error("expected:", etok)
		}
	}
}
