package parse

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var (
	tokens  = poLexer.Symbols()
	poLexer = lexer.MustSimple(
		[]lexer.SimpleRule{
			{Name: "WS", Pattern: "[\t\r ]+"},
			{Name: "STRING", Pattern: `"(\\"|[^"])*"`},
			{Name: "NL", Pattern: "\n"},
			{Name: "MSGCTXT", Pattern: "msgctxt"},
			{Name: "MSGID", Pattern: "msgid[^_]"},
			{Name: "MSGSTR", Pattern: `msgstr[^[]`},
			{Name: "PLURAL_MSGID", Pattern: "msgid_plural"},
			{Name: "PLURAL_MSGSTR", Pattern: `msgstr\[\d+\]`},
			{Name: "COMMENT", Pattern: "# *[^\n]*"},
			{Name: "FLAG_COMMENT", Pattern: "#,[^\n]*"},
			{Name: "EXTRACTED_COMMENT", Pattern: "#\\.[^\n]*"},
			{Name: "PREVIOUS_COMMENT", Pattern: "#\\|[^\n]*"},
			{Name: "REFERENCE_COMMENT", Pattern: "#:[^\n]*"},
		},
	)
	poParser = participle.MustBuild[poFile](
		participle.Lexer(poLexer),
		participle.Unquote("STRING"),
		participle.Elide(
			"WS",
			"NL",
			"COMMENT",
			"FLAG_COMMENT",
			"EXTRACTED_COMMENT",
			"PREVIOUS_COMMENT",
			"REFERENCE_COMMENT",
		),
	)
)

type poFile struct {
	Pos lexer.Position

	Entries []entry `@@*`
}

type entry struct {
	Tokens []lexer.Token

	Context     []string        `(MSGCTXT @STRING+)?`
	ID          []string        `MSGID @STRING+(`
	Str         []string        `MSGSTR @STRING+|`
	MsgidPlural []string        `(PLURAL_MSGID @STRING+`
	Plurals     []pluralEntries `@@*))`
}

type pluralEntries struct {
	ID  string   `@PLURAL_MSGSTR`
	Str []string `@STRING+`
}
