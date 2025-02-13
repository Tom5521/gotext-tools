package parse_test

import (
	"testing"

	"github.com/Tom5521/xgotext/pkg/po/parse"
)

func BenchmarkParse(b *testing.B) {
	input := []byte(`#, fuzzy
msgid ""
msgstr ""
"Project-Id-Version: \n"
"Report-Msgid-Bugs-To: \n"
"POT-Creation-Date: 2025-02-11 23:31: 15\n"
"PO-Revision-Date: \n"
"Last-Translator: \n"
"Language-Team: \n"
"Language: \n"
"MIME-Version: 1.0\n"
"Content-Type: text/plain; charset=CHARSET\n"
"Content-Transfer-Encoding: 8bit\n"
"Plural-Forms: nplurals=2; plural=(n != 1);\n"

msgid "Hello World!"
msgstr "¡Hola mundo!"

msgid "How are you?"
msgstr "¿Como estas?"

msgid "meow"
msgstr "miau"`)
	parser, err := parse.NewParserFromBytes(input, "test.po")
	if err != nil {
		b.Error(err)
		return
	}

	for range b.N {
		parser.Parse()
		if len(parser.Errors()) > 0 {
			b.Error(parser.Errors()[0])
		}
	}
}
