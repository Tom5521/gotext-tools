package parse

import (
	"fmt"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po/entry"
)

type Parser struct {
	lexer *Lexer
	file  *File
	name  string
}

func NewParser(input []rune, filename string) *Parser {
	p := &Parser{
		lexer: NewLexer(input),
		name:  filename,
		file:  new(File),
	}

	return p
}

func NewParserFromString(input, filename string) *Parser {
	return NewParser([]rune(input), filename)
}

func (p *Parser) genParseMap() map[Type]func(Token) (Node, error) {
	return map[Type]func(Token) (Node, error){
		COMMENT: p.Comment,
		MSGID:   p.Msgid,
		MSGSTR:  p.Msgstr,
		MSGCTXT: p.Msgctxt,
		// PluralMsgid:  p.PluralMsgid,
		// PluralMsgstr: p.PluralMsgstr,
	}
}

func (p *Parser) Parse() []error {
	var errs []error

	addErr := func(format string, a ...any) {
		errs = append(errs, fmt.Errorf(format, a...))
	}

	parseMap := p.genParseMap()

	for tok := p.lexer.NextToken(); tok.Type != EOF; tok = p.lexer.NextToken() {
		if tok.Type == ILLEGAL {
			addErr(
				"token at %s:%d is ILLEGAL",
				p.name,
				util.FindLine(p.lexer.input, tok.Pos),
			)
			continue
		}
		parse, ok := parseMap[tok.Type]
		if !ok {
			addErr("unknown token type at %s:%d", p.name, util.FindLine(p.lexer.input, tok.Pos))
			continue
		}

		node, err := parse(tok)
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
