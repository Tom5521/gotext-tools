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

	switch {
	case l.char == '#':
		tok.Type = COMMENT
		tok.Literal = l.readComment()
	case l.char == '"':
		tok.Type = STRING
		tok.Literal = l.readString()
	case l.isKeyword():
		tok.Literal = l.readKeyword()
		tok.Type = LookupIdent(tok.Literal)
	case l.char == 0:
		tok.Type = EOF
		tok.Literal = ""
	default:
		tok = Token{
			Type:    ILLEGAL,
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
