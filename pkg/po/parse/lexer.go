package parse

import (
	"unicode"
)

type Lexer struct {
	input []rune
	pos   int
	read  int
	char  rune
	prev  rune
}

func NewLexer(input []rune) *Lexer {
	l := &Lexer{
		input: input,
	}
	l.readChar()
	return l
}

func NewLexerFromString(input string) *Lexer {
	return NewLexer([]rune(input))
}

func (l *Lexer) readChar() {
	l.prev = l.char
	if l.read >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.read]
	}
	l.pos = l.read
	l.read++
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()
	tok := Token{
		Pos: l.pos,
	}

	switch l.char {
	case '#':
		tok.Type = COMMENT
		tok.Literal = l.readComment()
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
	case 0:
		tok.Type = EOF
		tok.Literal = ""
	default:
		switch {
		case unicode.IsLetter(l.char):
			tok.Literal = l.readKeyword()
			tok.Type = LookupIdent(tok.Literal)
		default:
			tok = Token{
				Type:    ILLEGAL,
				Literal: string(l.char),
			}
			l.readChar()
		}
	}

	return tok
}

func (l *Lexer) readPlural() string {
	pos := l.pos

	for (l.char != ']' || unicode.IsDigit(l.char)) && l.char != 0 {
		l.readChar()
	}
	l.readChar()

	return string(l.input[pos:l.pos])
}

func (l *Lexer) readDigit() string {
	pos := l.pos

	for unicode.IsDigit(l.char) {
		l.readChar()
	}

	return string(l.input[pos:l.pos])
}

func (l *Lexer) readString() string {
	pos := l.pos

	l.readChar() // Consume 1st quote.

	for l.char != '"' && l.char != '0' {
		l.readChar()
	}

	l.readChar() // Consume last quote.

	return string(l.input[pos:l.pos])
}

func (l *Lexer) readComment() string {
	pos := l.pos

	for l.char != '\n' && l.char != 0 {
		l.readChar()
	}

	return string(l.input[pos:l.pos])
}

func (l *Lexer) readKeyword() string {
	pos := l.pos
	for unicode.IsLetter(l.char) || l.char == '[' || l.char == '_' {
		if l.char == '[' {
			str := string(l.input[pos:l.pos])
			return str + l.readPlural()
		}
		l.readChar()
	}
	return string(l.input[pos:l.pos])
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.char) {
		l.readChar()
	}
}
