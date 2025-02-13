package ast_test

import (
	"testing"

	"github.com/Tom5521/xgotext/pkg/po/parse/ast"
)

var input = `# General Comment
#, flag comment
#: location_comment:123
#. extracted comment
msgctxt "testing!"
msgid "1st msgid!"
msgstr "1er msgid!"

msgid "I want an apple"
msgid_plural "I want some apples"
msgstr[0] "Quiero una manzana"
msgstr[1] "Quiero unas manzanas"

msgid ""
"hello"
"world"
""
msgstr ""
"hola"
"mundo"
""`

func BenchmarkParser(b *testing.B) {
	p := ast.NewParserFromString(input, "test.po")

	for range b.N {
		e := p.Parse()
		if len(e) > 0 {
			b.Error(e[0])
		}
	}
}

func BenchmarkNormalizer(b *testing.B) {
	p := ast.NewParserFromString(input, "test.po")
	n, e := p.Normalizer()
	if len(e) > 0 {
		b.Error(e[0])
		return
	}

	for range b.N {
		n.Normalize()

		if len(n.Errors()) > 0 {
			b.Error(n.Errors()[0])
		}
	}
}
