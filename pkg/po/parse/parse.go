package parse

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Tom5521/xgotext/internal/util"
)

func (p *Parser) readStringIdent() (string, error) {
	var b strings.Builder

	tok := p.tokens[p.position]

	next := p.token(p.position + 1)
	if next.Type != STRING {
		return "", fmt.Errorf(
			"expected STRING after %s declaration [%s:%d]",
			tok.Type,
			p.File.Name,
			util.FindLine(p.input, tok.Pos),
		)
	}

	str, err := strconv.Unquote(next.Literal)
	if err != nil {
		return "", err
	}

	b.WriteString(str)

	multiline, err := p.readMultilineStrings(2)
	if err != nil {
		return "", err
	}
	b.WriteString(multiline)

	return b.String(), nil
}

func (p *Parser) readMultilineStrings(off int) (string, error) {
	var b strings.Builder

	for _, tok := range p.tokens[p.position+off:] {
		if tok.Type != STRING {
			break
		}
		id, err := strconv.Unquote(tok.Literal)
		if err != nil {
			return "", err
		}
		b.WriteByte('\n')
		b.WriteString(id)
	}

	return b.String(), nil
}

func (p *Parser) comment() (Node, error) {
	tok := p.tokens[p.position]
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

func (p *Parser) msgid() (Node, error) {
	tok := p.tokens[p.position]
	msgid := Msgid{
		pos:     tok.Pos,
		literal: tok.Literal,
	}

	id, err := p.readStringIdent()
	if err != nil {
		return nil, err
	}

	msgid.ID = id

	return msgid, nil
}

func (p *Parser) msgstr() (Node, error) {
	tok := p.tokens[p.position]
	msgstr := Msgstr{
		pos:     tok.Pos,
		literal: tok.Literal,
	}

	str, err := p.readStringIdent()
	if err != nil {
		return nil, err
	}
	msgstr.Str = str

	return msgstr, nil
}

func (p *Parser) msgctxt() (Node, error) {
	tok := p.tokens[p.position]
	msgctxt := Msgctxt{
		pos:     tok.Pos,
		literal: tok.Literal,
	}

	ctx, err := p.readStringIdent()
	if err != nil {
		return nil, err
	}

	msgctxt.Context = ctx

	return msgctxt, nil
}

func (p *Parser) pluralMsgid() (Node, error) {
	tok := p.tokens[p.position]
	pmsgid := MsgidPlural{
		pos:     tok.Pos,
		literal: tok.Literal,
	}

	id, err := p.readStringIdent()
	if err != nil {
		return nil, err
	}

	pmsgid.Plural = id

	return pmsgid, nil
}

func (p *Parser) pluralMsgstr() (Node, error) {
	tok := p.tokens[p.position]
	pmsgstr := MsgstrPlural{
		pos:     tok.Pos,
		literal: tok.Literal,
	}

	var npluralID []rune
	for _, char := range tok.Literal[strings.Index(tok.Literal, "[")+1:] {
		if char == ']' {
			break
		}

		npluralID = append(npluralID, char)
	}

	id, err := strconv.Atoi(string(npluralID))
	if err != nil {
		return nil, err
	}
	pmsgstr.PluralID = id

	pmsgstr.Str, err = p.readStringIdent()
	if err != nil {
		return nil, err
	}

	return pmsgstr, nil
}
