package compile

import (
	"bytes"
	"fmt"
	"io"

	"github.com/Tom5521/gotext-tools/v2/internal/color"
	"github.com/Tom5521/gotext-tools/v2/internal/util"

	"github.com/alecthomas/participle/v2/lexer"
)

// CSSClassesHighlighting emulates CSS HighlightCSSProperties of https://www.gnu.org/software/gettext/manual/html_node/Style-rules.html
//
// Multiple classes are separated by spaces in the key field.
//
// Ex: CSSClassesHighlighting{"my-class1 my-class2":{Color: color.Red}}.
type CSSClassesHighlighting map[string]HighlightCSSProperties

// NOTE: I'll keep this only for internal documentation purpose.
//
/* struct {
	ID  *HighlightCSSProperties `css:"msgid"`
	Str *HighlightCSSProperties `css:"msgstr"`

	Fuzzy        *HighlightCSSProperties `css:"fuzzy"`
	Obsolete     *HighlightCSSProperties `css:"obsolete"`
	Translated   *HighlightCSSProperties `css:"translated"`
	Untranslated *HighlightCSSProperties `css:"untranslated"`

	Comment           *HighlightCSSProperties `css:"comment"`
	TranslatorComment *HighlightCSSProperties `css:"translator-comment"`
	ExtractedComment  *HighlightCSSProperties `css:"extracted-comment"`
	ReferenceComment  *HighlightCSSProperties `css:"reference-comment"`
	Reference         *HighlightCSSProperties `css:"reference"`
	FlagComment       *HighlightCSSProperties `css:"flag-comment"`
	Flag              *HighlightCSSProperties `css:"flag"`
	FuzzyFlag         *HighlightCSSProperties `css:"fuzzy-flag"`
	PreviousComment   *HighlightCSSProperties `css:"previous-comment"`
	Previous          *HighlightCSSProperties `css:"previous"`
	Keyword           *HighlightCSSProperties `css:"keyword"`
	String            *HighlightCSSProperties `css:"string"`

	Text                   *HighlightCSSProperties `css:"text"`
	EscapeSequence         *HighlightCSSProperties `css:"escape-sequence"`
	FormatDirective        *HighlightCSSProperties `css:"format-directive"`
	InvalidFormatDirective *HighlightCSSProperties `css:"invalid-format-directive"`
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

type TermColorer interface {
	Sprintf(string, ...any) string
	Sprint(...any) string
}

type HighlightCSSProperties struct {
	Color           TermColorer
	BackgroundColor TermColorer
	FontWeight      HighlightFontWeight
	FontStyle       HighlightFontStyle
	TextDecoration  HighlightTextDecoration
}

/*
IMPLEMENTED:
n = no
y = yes
~ = unfinished

msgid: n
msgstr: n
fuzzy: n
obsolete: n
translated: n
untranslated: n
comment: n
translator-comment: n
extracted-comment: n
reference-comment: n
reference: n
flag-comment: n
flag: n
fuzzy-flag: n
previous-comment: n
previous: n
keyword: y
string: y
text: y
escape-sequence: n
format-directive: n
invalid-format-directive: n
"x y": ~
*/
var TestHighlight = CSSClassesHighlighting{
	"msgid":  {Color: color.Blue, FontWeight: FontWeightBold},
	"msgstr": {Color: color.Green, FontStyle: FontStyleItalic},

	"fuzzy":        {Color: color.Yellow, FontStyle: FontStyleItalic},
	"obsolete":     {Color: color.Green, FontStyle: FontStyleItalic},
	"translated":   {Color: color.Green, FontWeight: FontWeightBold},
	"untranslated": {Color: color.Red, FontWeight: FontWeightNormal},

	"comment":            {Color: color.Cyan, FontStyle: FontStyleItalic},
	"translator-comment": {Color: color.Green},
	"extracted-comment":  {Color: color.Green, FontWeight: FontWeightBold},
	"reference-comment":  {Color: color.Magenta, FontStyle: FontStyleItalic},
	"reference":          {Color: color.Magenta, FontWeight: FontWeightBold},
	"flag-comment":       {Color: color.Yellow, TextDecoration: TextDecorationUnderline},
	"flag":               {TextDecoration: TextDecorationUnderline},
	"fuzzy-flag":         {TextDecoration: TextDecorationNone},
	"previous-comment":   {Color: color.Green, FontStyle: FontStyleItalic},
	"previous":           {Color: color.Green},
	"keyword":            {Color: color.Blue, FontWeight: FontWeightBold},
	"string":             {Color: color.Green},

	"text": {
		Color: color.Red,
		// BackgroundColor: color.BgWhite,
		FontWeight: FontWeightBold,
	},
	"escape-sequence":  {Color: color.Red, FontWeight: FontWeightBold},
	"format-directive": {FontWeight: FontWeightBold},
	"invalid-format-directive": {
		BackgroundColor: color.Red,
		Color:           color.White,
		FontWeight:      FontWeightBold,
	},

	"msgid text":        {Color: color.Blue},
	"msgstr text":       {Color: color.Green},
	"fuzzy msgstr text": {Color: color.Red},
}

var DefaultHighlight = CSSClassesHighlighting{
	"translator-comment": {Color: color.Green},
	"obsolete":           {Color: color.Green},
	"extracted-comment":  {Color: color.Green, FontWeight: FontWeightBold},
	"flag":               {TextDecoration: TextDecorationUnderline},
	"fuzzy-flag":         {TextDecoration: TextDecorationNone},
	"text":               {Color: color.Magenta},
	"msgstr text":        {Color: color.Blue},
	"fuzzy msgstr text":  {Color: color.Red},
	"format-directive":   {FontWeight: FontWeightBold},
	"invalid-format-directive": {
		BackgroundColor: color.Red,
		Color:           color.White,
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

// TODO: Finish this.
func highlight(cfg CSSClassesHighlighting, lex lexer.Lexer) ([]byte, error) {
	tokens, err := lexer.ConsumeAll(lex)
	if err != nil {
		return nil, fmt.Errorf("error consuming lexer tokens: %w", err)
	}

	// Rebuild file.
	var builder bytes.Buffer
	for _, t := range tokens {
		builder.WriteString(t.Value)
	}

	return builder.Bytes(), nil
}
