package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po/parse/token"
)

func (p *Parser) readStringIdent() (string, error) {
	var b strings.Builder

	current := p.tokens[p.position]

	var lines []string
	for _, tok := range p.tokens[p.position+1:] {
		if tok.Type != token.STRING {
			break
		}
		id, err := strconv.Unquote(tok.Literal)
		if err != nil {
			return "", err
		}

		lines = append(lines, id)
	}

	if len(lines) == 0 {
		return "", fmt.Errorf(
			"expected STRING after %s declaration [%s:%d]",
			current.Type,
			p.File.Name,
			util.FindLine(p.input, current.Pos),
		)
	}

	for i, line := range lines {
		if i == len(lines)-1 {
			fmt.Fprint(&b, line)
			continue
		}
		fmt.Fprintln(&b, line)
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
