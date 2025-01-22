package parse

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Tom5521/xgotext/internal/util"
)

func (p *Parser) readStringIdent(tok Token) (string, error) {
	var b strings.Builder

	nextTok := p.lexer.NextToken()
	if nextTok.Type != STRING {
		return "", fmt.Errorf(
			"expected STRING after %s declaration [%s:%d]",
			tok.Type,
			p.name,
			util.FindLine(p.lexer.input, tok.Pos),
		)
	}

	str, err := strconv.Unquote(nextTok.Literal)
	if err != nil {
		return "", err
	}

	b.WriteString(str)

	multiline, err := p.readMultilineStrings()
	if err != nil {
		return "", err
	}
	b.WriteString(multiline)

	return b.String(), nil
}

func (p *Parser) readMultilineStrings() (string, error) {
	var b strings.Builder

	for t := p.lexer.NextToken(); t.Type == STRING; t = p.lexer.NextToken() {
		id, err := strconv.Unquote(t.Literal)
		if err != nil {
			return "", err
		}
		b.WriteByte('\n')
		b.WriteString(id)
	}

	return b.String(), nil
}

func (p *Parser) Comment(tok Token) (Node, error) {
	if len(tok.Literal) == 1 {
		return GeneralComment{
			pos:     tok.Pos,
			literal: tok.Literal,
		}, nil
	}

	parts := strings.Fields(tok.Literal)

	switch parts[0] {
	case "#:":
		info := strings.SplitN(strings.Join(parts[1:], ""), ":", 2)
		line, err := strconv.Atoi(info[1])
		if err != nil {
			return nil, err
		}

		return LocationComment{
			pos:     tok.Pos,
			literal: tok.Literal,
			File:    info[0],
			Line:    line,
		}, nil
	case "#,":
		return FlagComment{
			pos:     tok.Pos,
			literal: tok.Literal,
			Flag:    strings.Join(parts[1:], " "),
		}, nil
	default:
		return GeneralComment{
			pos:     tok.Pos,
			literal: tok.Literal,
			Text:    strings.Join(parts[1:], " "),
		}, nil
	}
}

func (p *Parser) Msgid(tok Token) (Node, error) {
	msgid := Msgid{
		pos:     tok.Pos,
		literal: tok.Literal,
	}

	id, err := p.readStringIdent(tok)
	if err != nil {
		return nil, err
	}

	msgid.ID = id

	return msgid, nil
}

func (p *Parser) Msgstr(tok Token) (Node, error) {
	msgstr := Msgstr{
		pos:     tok.Pos,
		literal: tok.Literal,
	}

	str, err := p.readStringIdent(tok)
	if err != nil {
		return nil, err
	}
	msgstr.Str = str

	return msgstr, nil
}

func (p *Parser) Msgctxt(tok Token) (Node, error) {
	msgctxt := &Msgctxt{
		pos:     tok.Pos,
		literal: tok.Literal,
	}

	ctx, err := p.readStringIdent(tok)
	if err != nil {
		return nil, err
	}

	msgctxt.Context = ctx

	return msgctxt, nil
}

// TODO: Finish this.
// func p.PluralMsgid(tok Token) (Node, error)
// func p.PluralMsgstr(tok Token) (Node, error)
