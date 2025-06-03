package util

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
		MsgidPlural []string        `| (MsgidPlural @String+`
		Plurals     []pluralEntries `@@*))`
	}

	pluralEntries struct {
		ID  int      `Msgstr LB @Integer RB`
		Str []string `@String+`
	}
)

func SearchSymbol(t lexer.TokenType) string {
	for k, tt := range PoSymbols {
		if tt == t {
			return k
		}
	}
	return ""
}

var (
	PoSymbols = PoLexer.Symbols()
	PoRules   = []lexer.SimpleRule{
		{Name: "WS", Pattern: `\s+`},
		{Name: "Integer", Pattern: `\d+`},
		{Name: "LB", Pattern: `\[`},
		{Name: "RB", Pattern: `\]`},
		{Name: "String", Pattern: `"(\\"|[^"])*"`},
		{Name: "Msgctxt", Pattern: "msgctxt"},
		{Name: "MsgidPlural", Pattern: "msgid_plural"},
		{Name: "Msgid", Pattern: "msgid"},
		{Name: "Msgstr", Pattern: "msgstr"},
		{Name: "Comment", Pattern: "#[^\n]*"},
	}
	PoLexer  = lexer.MustSimple(PoRules)
	PoParser = participle.MustBuild[poFile](
		participle.Lexer(PoLexer),
		participle.Unquote("String"),
		participle.Elide(
			"WS",
			"Comment",
		),
	)
)
