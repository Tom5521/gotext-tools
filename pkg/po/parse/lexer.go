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
	var tok Token
	l.skipWhitespace()

	switch l.char {
	case '#':
		tok.Type = COMMENT
		tok.Literal = l.readComment()
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
		l.readChar() // consume closing quote
	case 0:
		tok.Type = EOF
		tok.Literal = ""
	default:
		if unicode.IsLetter(l.char) {
			tok.Literal = l.readKeyword()
			t, ok := keywords[tok.Literal]
			if ok {
				tok.Type = t
			} else {
				tok.Type = ILLEGAL
			}
			return tok
		} else {
			tok = Token{
				Type:    ILLEGAL,
				Literal: string(l.char),
			}
			l.readChar()
		}
	}

	return tok
}

func (l *Lexer) readString() string {
	l.readChar() // skip opening quote
	pos := l.pos

	for l.char != '"' && l.char != 0 {
		l.readChar()
	}

	return string(l.input[pos:l.pos])
}

func (l *Lexer) readComment() string {
	l.readChar() // skip #
	pos := l.pos

	for l.char != '\n' && l.char != 0 {
		l.readChar()
	}

	return string(l.input[pos:l.pos])
}

func (l *Lexer) readKeyword() string {
	pos := l.pos
	for unicode.IsLetter(l.char) || unicode.IsDigit(l.char) {
		l.readChar()
	}
	return string(l.input[pos:l.pos])
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.char) {
		l.readChar()
	}
}
