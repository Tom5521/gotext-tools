package parse

import (
	"fmt"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po/entry"
)

type Parser struct {
	input    []rune
	tokens   []Token
	position int
	file     *File
	name     string
}

func (p *Parser) collectTokens(l *Lexer) {
	tok := l.NextToken()
	for tok.Type != EOF {
		p.tokens = append(p.tokens, tok)
		tok = l.NextToken()
	}
}

func NewParser(input []rune, filename string) *Parser {
	p := &Parser{
		input: input,
		name:  filename,
		file:  new(File),
	}

	p.collectTokens(NewLexer(input))

	return p
}

func NewParserFromString(input, filename string) *Parser {
	return NewParser([]rune(input), filename)
}

func (p *Parser) genParseMap() map[Type]func() (Node, error) {
	return map[Type]func() (Node, error){
		COMMENT:      p.comment,
		MSGID:        p.msgid,
		MSGSTR:       p.msgstr,
		MSGCTXT:      p.msgctxt,
		PluralMsgid:  p.pluralMsgid,
		PluralMsgstr: p.pluralMsgstr,
	}
}

func (p *Parser) token(i int) Token {
	if i < 0 || i >= len(p.tokens) {
		return Token{Type: EOF}
	}

	return p.tokens[i]
}

func (p *Parser) Parse() []error {
	var errs []error

	parseMap := p.genParseMap()

	for i, tok := range p.tokens {
		p.position = i
		var node Node
		var err error
		switch tok.Type {
		case ILLEGAL:
			err = fmt.Errorf(
				"token at %s:%d is illegal",
				p.name,
				util.FindLine(p.input, tok.Pos),
			)
		case MSGID, MSGSTR, MSGCTXT, PluralMsgid, PluralMsgstr, COMMENT:
			parser := parseMap[tok.Type]
			node, err = parser()
		case STRING:
			continue
		default:
			err = fmt.Errorf(
				"unknown token type at %s:%d",
				p.name,
				util.FindLine(p.input, tok.Pos),
			)
		}

		if err != nil {
			errs = append(errs, err)
			continue
		}

		p.file.Nodes = append(p.file.Nodes, node)
	}

	return errs
}

func (p *Parser) Nodes() []Node {
	return p.file.Nodes
}

func (p *Parser) Translations() []entry.Translation { return nil }
