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

var idTokensMap = map[lexer.TokenType]bool{
	util.PoSymbols["Msgid"]:   false,
	util.PoSymbols["Msgctxt"]: false,
	util.PoSymbols["Plural"]:  false,
}

var strTokensMap = map[lexer.TokenType]bool{
	util.PoSymbols["Msgstr"]: false,
	util.PoSymbols["RB"]:     false,
}

func HighlightOutput(cfg *HighlightConfig, name string, input io.Reader) ([]byte, error) {
	lex, err := util.PoLexer.Lex(name, input)
	if err != nil {
		return nil, err
	}

	tokens, err := lexer.ConsumeAll(lex)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(tokens); i++ {
		ttype := tokens[i].Type

		if ttype == util.PoSymbols["Comment"] {
			tokens[i].Value = cfg.Comment.Render(tokens[i].Value)
			continue
		}

		if _, ok := idTokensMap[ttype]; ok {
			i += colorStrings(tokens, i+1, cfg.ID, cfg.Comment)
		} else if _, ok := strTokensMap[ttype]; ok {
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
