package compile

import (
	"bytes"
	"fmt"
	"io"
	"regexp"

	"github.com/Tom5521/gotext-tools/v2/internal/util"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/gookit/color"
)

type HighlightConfig struct {
	ID, Str, Comment color.Color
}

var DefaultHighlight = &HighlightConfig{color.Magenta, color.Blue, color.Green}

func Highlight(cfg *HighlightConfig, name, input string) ([]byte, error) {
	lex, err := util.PoLexer.LexString(name, input)
	if err != nil {
		return nil, fmt.Errorf("error lexing string: %w", err)
	}
	return highlight(cfg, lex)
}

func HighlightFromBytes(cfg *HighlightConfig, name string, input []byte) ([]byte, error) {
	lex, err := util.PoLexer.LexString(name, string(input))
	if err != nil {
		return nil, fmt.Errorf("error lexing bytes: %w", err)
	}
	return highlight(cfg, lex)
}

func HighlightFromReader(cfg *HighlightConfig, name string, input io.Reader) ([]byte, error) {
	lex, err := util.PoLexer.Lex(name, input)
	if err != nil {
		return nil, fmt.Errorf("error lexing reader: %w", err)
	}

	return highlight(cfg, lex)
}

var idTokensMap = map[lexer.TokenType]struct{}{
	util.PoSymbols["Msgid"]:   {},
	util.PoSymbols["Msgctxt"]: {},
	util.PoSymbols["Plural"]:  {},
}

var strTokensMap = map[lexer.TokenType]struct{}{
	util.PoSymbols["Msgstr"]: {},
	util.PoSymbols["RB"]:     {},
}

func highlight(cfg *HighlightConfig, lex lexer.Lexer) ([]byte, error) {
	tokens, err := lexer.ConsumeAll(lex)
	if err != nil {
		return nil, fmt.Errorf("error consuming lexer tokens: %w", err)
	}

	for i := 0; i < len(tokens); i++ {
		ttype := tokens[i].Type

		if ttype == util.PoSymbols["Comment"] {
			tokens[i].Value = cfg.Comment.Render(tokens[i].Value)
			continue
		}

		if _, ok := idTokensMap[ttype]; ok {
			i += colorStrings(tokens, i+1, cfg.ID, cfg.Comment)
		} else if _, ok = strTokensMap[ttype]; ok {
			i += colorStrings(tokens, i+1, cfg.Str, cfg.Comment)
		}
	}

	// Rebuild file.
	var builder bytes.Buffer
	for _, t := range tokens {
		builder.WriteString(t.Value)
	}

	return builder.Bytes(), nil
}

var strRegex = regexp.MustCompile(`"(.*)"`)

func colorStrings(tokens []lexer.Token, offset int, unq, comment color.Color) int {
	var mod int

	for i := offset; i < len(tokens); i++ {
		t := tokens[i]

		switch t.Type {
		case util.PoSymbols["WS"]:
			continue
		case util.PoSymbols["Comment"]:
			t.Value = comment.Render(t.Value)
		case util.PoSymbols["String"]:
			// NOTE:
			// The lexer guarantees that tokens of type "String" are always properly quoted.
			// Therefore, it's safe to access the first capture group without additional checks.
			unquoted := strRegex.FindStringSubmatch(t.Value)[1]
			t.Value = fmt.Sprintf(`"%s"`, unq.Render(unquoted))
		default:
			return mod
		}

		mod++
		tokens[i] = t
	}

	return mod
}
