package compile

import (
	"bytes"
	"fmt"
	"io"
	"regexp"

	"github.com/Tom5521/gotext-tools/v2/internal/util"

	"github.com/alecthomas/participle/v2/lexer"
)

// HighlightColor represents an ANSI text colorizer.
// I leave it as an interface so that anyone can choose
// the library they want to use to color the text.
//
// The default implementation is a very basic implementation of a colorizer.
//
// This is intended to be used alongside CSS,
// but in practice it can be used as you wish.
//
// WARNING: If the colorizer goes beyond ANSI, you're on your luck.
type HighlightColor interface {
	Sprint(...any) string
	Sprintf(string, ...any) string
}

type hcolor = HighlightColor

// CSSClassesHighlighting emulates CSS properties of https://www.gnu.org/software/gettext/manual/html_node/Style-rules.html
//
// Multiple classes are separated by spaces in the key field.
//
// Ex: CSSClassesHighlighting{"my-class1 my-class2":{Color: color.Red}}
type CSSClassesHighlighting map[string]properties

// NOTE: I'll keep this only for internal documentation purpose.
//
/* struct {
	ID  *properties `css:"msgid"`
	Str *properties `css:"msgstr"`

	Fuzzy        *properties `css:"fuzzy"`
	Obsolete     *properties `css:"obsolete"`
	Translated   *properties `css:"translated"`
	Untranslated *properties `css:"untranslated"`

	Comment           *properties `css:"comment"`
	TranslatorComment *properties `css:"translator-comment"`
	ExtractedComment  *properties `css:"extracted-comment"`
	ReferenceComment  *properties `css:"reference-comment"`
	Reference         *properties `css:"reference"`
	FlagComment       *properties `css:"flag-comment"`
	Flag              *properties `css:"flag"`
	FuzzyFlag         *properties `css:"fuzzy-flag"`
	PreviousComment   *properties `css:"previous-comment"`
	Previous          *properties `css:"previous"`
	Keyword           *properties `css:"keyword"`
	String            *properties `css:"string"`

	Text                   *properties `css:"text"`
	EscapeSequence         *properties `css:"escape-sequence"`
	FormatDirective        *properties `css:"format-directive"`
	InvalidFormatDirective *properties `css:"invalid-format-directive"`
} */

type HighlightFontStyle int

const (
	FontStyleItalic HighlightFontStyle = 1 + iota
	FontStyleOblique
)

type HighlightTextDecoration int

const (
	TextDecorationNone HighlightTextDecoration = iota
	TextDecorationUnderline
)

type HighlightFontWeight int

const (
	FontWeightNormal = iota
	FontWeightBold
)

type HighlightCSSProperties struct {
	Color           hcolor
	BackgroundColor hcolor
	FontWeight      HighlightFontWeight
	FontStyle       HighlightFontStyle
	TextDecoration  HighlightTextDecoration
}

type properties = HighlightCSSProperties

var DefaultHighlight = CSSClassesHighlighting{
	"translator-comment": {Color: util.Green},
	"obsolete":           {Color: util.Green},
	"extracted-comment":  {Color: util.Green, FontWeight: FontWeightBold},
	"flag":               {TextDecoration: TextDecorationUnderline},
	"fuzzy-flag":         {TextDecoration: TextDecorationNone},
	"text":               {Color: util.Magenta},
	"msgstr text":        {Color: util.Blue},
	"fuzzy msgstr text":  {Color: util.Red},
	"format-directive":   {FontWeight: FontWeightBold},
	"invalid-format-directive": {
		BackgroundColor: util.Red,
		Color:           util.White,
		FontWeight:      FontWeightBold,
	},
}

func Highlight(cfg CSSClassesHighlighting, name, input string) ([]byte, error) {
	lex, err := util.PoLexer.LexString(name, input)
	if err != nil {
		return nil, fmt.Errorf("error lexing string: %w", err)
	}
	return highlight(cfg, lex)
}

func HighlightFromBytes(cfg CSSClassesHighlighting, name string, input []byte) ([]byte, error) {
	lex, err := util.PoLexer.LexString(name, string(input))
	if err != nil {
		return nil, fmt.Errorf("error lexing bytes: %w", err)
	}
	return highlight(cfg, lex)
}

func HighlightFromReader(
	cfg CSSClassesHighlighting,
	name string,
	input io.Reader,
) ([]byte, error) {
	lex, err := util.PoLexer.Lex(name, input)
	if err != nil {
		return nil, fmt.Errorf("error lexing reader: %w", err)
	}

	return highlight(cfg, lex)
}

/* var idTokensMap = map[lexer.TokenType]struct{}{
	util.PoSymbols["Msgid"]:   {},
	util.PoSymbols["Msgctxt"]: {},
	util.PoSymbols["Plural"]:  {},
}

var strTokensMap = map[lexer.TokenType]struct{}{
	util.PoSymbols["Msgstr"]: {},
	util.PoSymbols["RB"]:     {},
} */

// TODO: Finish this.
func highlight(cfg CSSClassesHighlighting, lex lexer.Lexer) ([]byte, error) {
	tokens, err := lexer.ConsumeAll(lex)
	if err != nil {
		return nil, fmt.Errorf("error consuming lexer tokens: %w", err)
	}

	/* 	for i := 0; i < len(tokens); i++ {
		ttype := tokens[i].Type

		if ttype == util.PoSymbols["Comment"] {
			tokens[i].Value = cfg.Comment.Color.Render(tokens[i].Value)
			continue
		}

		if _, ok := idTokensMap[ttype]; ok {
			if cfg.ID != nil {
				i += colorStrings(tokens, i+1, cfg.ID.Color, cfg.Comment.Color)
			}
		} else if _, ok = strTokensMap[ttype]; ok {
			if cfg.Str != nil {
				i += colorStrings(tokens, i+1, cfg.Str.Color, cfg.Comment.Color)
			}
		}
	} */

	// Rebuild file.
	var builder bytes.Buffer
	for _, t := range tokens {
		builder.WriteString(t.Value)
	}

	return builder.Bytes(), nil
}

var strRegex = regexp.MustCompile(`"(.*)"`)

func colorStrings(tokens []lexer.Token, offset int, unq, comment hcolor) int {
	var mod int

	for i := offset; i < len(tokens); i++ {
		t := tokens[i]

		switch t.Type {
		case util.PoSymbols["WS"]:
			continue
		case util.PoSymbols["Comment"]:
			t.Value = comment.Sprint(t.Value)
		case util.PoSymbols["String"]:
			// NOTE:
			// The lexer guarantees that tokens of type "String" are always properly quoted.
			// Therefore, it's safe to access the first capture group without additional checks.
			unquoted := strRegex.FindStringSubmatch(t.Value)[1]
			t.Value = fmt.Sprintf(`"%s"`, unq.Sprint(unquoted))
		default:
			return mod
		}

		mod++
		tokens[i] = t
	}

	return mod
}
