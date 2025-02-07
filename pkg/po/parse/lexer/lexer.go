package lexer

import (
	"unicode"

	"github.com/Tom5521/xgotext/pkg/po/parse/token"
)

type Lexer struct {
	input []byte
	pos   int
	read  int
	char  rune
	prev  rune
}

func New(input []byte) *Lexer {
	l := &Lexer{
		input: input,
	}
	l.readChar()
	return l
}

func NewFromString(input string) *Lexer {
	return New([]byte(input))
}

func (l *Lexer) readChar() {
	l.prev = l.char
	if l.read >= len(l.input) {
		l.char = 0
	} else {
		l.char = rune(l.input[l.read])
	}
	l.pos = l.read
	l.read++
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()
	tok := token.Token{
		Pos: l.pos,
	}

	switch {
	case l.char == '#':
		tok.Type = token.COMMENT
		tok.Literal = l.readComment()
	case l.char == '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case l.isKeyword():
		tok.Literal = l.readKeyword()
		tok.Type = token.LookupIdent(tok.Literal)
	case l.char == 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		tok = token.Token{
			Type:    token.ILLEGAL,
			Literal: string(l.char),
		}
		l.readChar()
	}

	return tok
}

func (l *Lexer) isKeyword() bool {
	return unicode.IsLetter(l.char)
}

func (l *Lexer) readKeyword() string {
	pos := l.pos

	for l.char != 0 && (unicode.IsLetter(l.char) || l.char == '_' || l.char == '[' || l.char == ']' || unicode.IsDigit(l.char)) {
		l.readChar()
	}
	return string(l.input[pos:l.pos])
}

func (l *Lexer) readString() string {
	pos := l.pos
	l.readChar() // Consume opening quote

	for l.char != '"' && l.char != 0 {
		if l.char == '\\' {
			l.readChar() // Skip escape character
			l.readChar() // Skip escaped character
			continue
		}
		l.readChar()
	}

	l.readChar() // Consume closing quote
	return string(l.input[pos:l.pos])
}

func (l *Lexer) readComment() string {
	pos := l.pos

	// Handle different comment types: translator, extracted, reference, flags, previous
	for l.char != '\n' && l.char != 0 {
		l.readChar()
	}

	return string(l.input[pos:l.pos])
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.char) {
		l.readChar()
	}
}
