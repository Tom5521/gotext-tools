package ast

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po/parse/token"
)

func (p *Tokenizer) readStringIdent() (string, error) {
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
			p.name,
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

var (
	locationRegex  = regexp.MustCompile(`#:\s?(.*)`)
	generalRegex   = regexp.MustCompile(`#\s?(.*)`)
	extractedRegex = regexp.MustCompile(`#\.\s?(.*)`)
	flagRegex      = regexp.MustCompile(`#,\s?(.*)`)
	previousRegex  = regexp.MustCompile(`#\|\s?(.*)`)
)

func (p *Tokenizer) comment() (Node, error) {
	tok := p.tokens[p.position]

	switch {
	case locationRegex.MatchString(tok.Literal):
		matches := locationRegex.FindStringSubmatch(tok.Literal)
		parts := strings.SplitN(matches[1], ":", 2)
		line := -1
		var err error
		if parts[1] != "" {
			line, err = strconv.Atoi(parts[1])
			if err != nil {
				return nil, err
			}
		}

		return LocationComment{
			pos:  tok.Pos,
			Line: line,
			File: parts[0],
		}, nil
	case extractedRegex.MatchString(tok.Literal):
		return ExtractedComment{
			pos:  tok.Pos,
			Text: extractedRegex.FindStringSubmatch(tok.Literal)[1],
		}, nil
	case flagRegex.MatchString(tok.Literal):
		return FlagComment{
			pos:  tok.Pos,
			Flag: flagRegex.FindStringSubmatch(tok.Literal)[1],
		}, nil
	case previousRegex.MatchString(tok.Literal):
		return PreviousComment{
			pos:  tok.Pos,
			Text: previousRegex.FindStringSubmatch(tok.Literal)[1],
		}, nil
	default:
		return GeneralComment{
			pos:  tok.Pos,
			Text: generalRegex.FindStringSubmatch(tok.Literal)[1],
		}, nil
	}
}

func (p *Tokenizer) msgid() (Node, error) {
	tok := p.tokens[p.position]
	msgid := Msgid{
		pos: tok.Pos,
	}

	id, err := p.readStringIdent()
	if err != nil {
		return nil, err
	}

	msgid.ID = id

	return msgid, nil
}

func (p *Tokenizer) msgstr() (Node, error) {
	tok := p.tokens[p.position]
	msgstr := Msgstr{
		pos: tok.Pos,
	}

	str, err := p.readStringIdent()
	if err != nil {
		return nil, err
	}
	msgstr.Str = str

	return msgstr, nil
}

func (p *Tokenizer) msgctxt() (Node, error) {
	tok := p.tokens[p.position]
	msgctxt := Msgctxt{
		pos: tok.Pos,
	}

	ctx, err := p.readStringIdent()
	if err != nil {
		return nil, err
	}

	msgctxt.Context = ctx

	return msgctxt, nil
}

func (p *Tokenizer) pluralMsgid() (Node, error) {
	tok := p.tokens[p.position]
	pmsgid := MsgidPlural{
		pos: tok.Pos,
	}

	id, err := p.readStringIdent()
	if err != nil {
		return nil, err
	}

	pmsgid.Plural = id

	return pmsgid, nil
}

var pluralRegex = regexp.MustCompile(`msgstr\[(\d*)\]`)

func (p *Tokenizer) pluralMsgstr() (Node, error) {
	tok := p.tokens[p.position]
	pmsgstr := MsgstrPlural{
		pos: tok.Pos,
	}

	npluralID := pluralRegex.FindStringSubmatch(tok.Literal)[1]

	id, err := strconv.Atoi(npluralID)
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
