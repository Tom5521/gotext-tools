package compile

import (
	"bytes"
	"fmt"
	"io"
	"regexp"

	"github.com/Tom5521/gotext-tools/v2/internal/slices"
	"github.com/Tom5521/gotext-tools/v2/internal/util"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/gookit/color"
)

type HighlightConfig struct {
	ID, Str, Comment color.Color
}

var DefaultHighligh = &HighlightConfig{
	color.Magenta, color.Blue, color.Green,
}

func HighlighOutput(cfg *HighlightConfig, name string, input io.Reader) ([]byte, error) {
	lex, err := util.PoLexer.Lex(name, input)
	if err != nil {
		return nil, err
	}

	tokens, err := lexer.ConsumeAll(lex)
	if err != nil {
		return nil, err
	}

	idTokens := []lexer.TokenType{
		util.PoSymbols["Msgid"],
		util.PoSymbols["Msgctxt"],
		util.PoSymbols["Plural"],
	}

	strTokens := []lexer.TokenType{
		util.PoSymbols["Msgstr"],
		util.PoSymbols["RB"],
	}

	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type == util.PoSymbols["Comment"] {
			tokens[i].Value = cfg.Comment.Render(tokens[i].Value)
		}
		switch {
		case slices.Contains(idTokens, tokens[i].Type):
			i += colorStrings(tokens, i+1, cfg.ID, cfg.Comment)
		case slices.Contains(strTokens, tokens[i].Type):
			i += colorStrings(tokens, i+1, cfg.Str, cfg.Comment)
		}
	}

	var builder bytes.Buffer
	for _, t := range tokens {
		builder.WriteString(t.Value)
	}
	return builder.Bytes(), nil
}

func colorStrings(tokens []lexer.Token, offset int, unq, comment color.Color) int {
	var mod int
	for i := offset; i < len(tokens); i++ {
		t := tokens[i]

		regex := regexp.MustCompile(`"(.*)"`)
		var unquoted string

		if t.Type == util.PoSymbols["WS"] {
			goto finish
		}

		if t.Type == util.PoSymbols["Comment"] {
			t.Value = comment.Render(t.Value)
			goto finish
		}
		if t.Type != util.PoSymbols["String"] {
			break
		}

		unquoted = regex.FindStringSubmatch(t.Value)[1]
		t.Value = fmt.Sprintf(`"%s"`, unq.Render(unquoted))

	finish:
		mod++
		tokens[i] = t
	}

	return mod
}
