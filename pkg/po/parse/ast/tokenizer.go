package ast

import (
	"errors"
	"fmt"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po/parse/lexer"
	"github.com/Tom5521/xgotext/pkg/po/parse/token"
)

type Tokenizer struct {
	input    []byte
	tokens   []token.Token
	position int

	nodes  []Node
	errors []error
	name   string
}

func (p *Tokenizer) collectTokens(l *lexer.Lexer) {
	tok := l.NextToken()
	for tok.Type != token.EOF {
		p.tokens = append(p.tokens, tok)
		tok = l.NextToken()
	}
}

func NewTokenizer(input []byte, filename string) *Tokenizer {
	p := &Tokenizer{
		input: input,
		name:  filename,
	}

	p.collectTokens(lexer.New(input))

	return p
}

func NewTokenizerFromLexer(l *lexer.Lexer, input []byte, name string) *Tokenizer {
	p := &Tokenizer{
		input: input,
		name:  name,
	}
	p.collectTokens(l)

	return p
}

func NewTokenizerFromString(input, filename string) *Tokenizer {
	return NewTokenizer([]byte(input), filename)
}

type parserFunc = func() (Node, error)

func (p *Tokenizer) genParseMap() map[token.Type]parserFunc {
	return map[token.Type]parserFunc{
		token.COMMENT:      p.comment,
		token.MSGID:        p.msgid,
		token.MSGSTR:       p.msgstr,
		token.MSGCTXT:      p.msgctxt,
		token.PluralMsgid:  p.pluralMsgid,
		token.PluralMsgstr: p.pluralMsgstr,
	}
}

func (p *Tokenizer) Normalizer() (*ASTBuilder, []error) {
	p.Tokenize()
	return NewASTBuilder(p.name, p.input, p.nodes), p.Errors()
}

func (p *Tokenizer) Tokenize() {
	p.nodes = nil
	p.errors = nil

	parseMap := p.genParseMap()

	var tok token.Token
	for p.position, tok = range p.tokens {
		if len(p.errors) > 3 {
			p.errors = append(p.errors, errors.New("too many errors"))
			break
		}
		var node Node
		var err error
		switch tok.Type {
		case token.ILLEGAL:
			err = fmt.Errorf(
				"token at %s:%d is illegal",
				p.name,
				util.FindLine(p.input, tok.Pos),
			)
		case token.MSGID,
			token.MSGSTR,
			token.MSGCTXT,
			token.PluralMsgid,
			token.PluralMsgstr,
			token.COMMENT:
			parser := parseMap[tok.Type]
			node, err = parser()
		case token.STRING:
			continue
		default:
			err = fmt.Errorf(
				"unknown token type at %s:%d",
				p.name,
				util.FindLine(p.input, tok.Pos),
			)
		}

		if err != nil {
			p.errors = append(p.errors, err)
			continue
		}

		p.nodes = append(p.nodes, node)
	}
}

func (p Tokenizer) Nodes() []Node {
	return p.nodes
}

func (p Tokenizer) Errors() []error {
	return p.errors
}
