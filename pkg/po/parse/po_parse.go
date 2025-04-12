package parse

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type (
	poFile struct {
		Tokens  []lexer.Token
		Entries []entry `@@*`
	}

	entry struct {
		Tokens []lexer.Token

		Context     []string        `(Msgctxt @String+)?`
		ID          []string        `Msgid @String+`
		Str         []string        `(Msgstr @String+`
		MsgidPlural []string        `| (Msgid Plural @String+`
		Plurals     []pluralEntries `@@*))`
	}

	pluralEntries struct {
		ID  int      `Msgstr LB @Integer RB`
		Str []string `@String+`
	}
)

var (
	symbols = poLexer.Symbols()
	poLexer = lexer.MustSimple(
		[]lexer.SimpleRule{
			{Name: "WS", Pattern: `\s+`},
			{Name: "Integer", Pattern: `\d+`},
			{Name: "LB", Pattern: `\[`},
			{Name: "RB", Pattern: `\]`},
			{Name: "String", Pattern: `"(\\"|[^"])*"`},
			{Name: "Msgctxt", Pattern: "msgctxt"},
			{Name: "Msgid", Pattern: "msgid"},
			{Name: "Msgstr", Pattern: "msgstr"},
			{Name: "Plural", Pattern: "_plural"},
			{Name: "Comment", Pattern: "#[^\n]*"},
		},
	)
	poParser = participle.MustBuild[poFile](
		participle.Lexer(poLexer),
		participle.Unquote("String"),
		participle.Elide(
			"WS",
			"Comment",
		),
	)
)
