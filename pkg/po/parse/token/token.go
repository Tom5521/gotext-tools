// Package token defines the types and constants used
// to represent tokens in the lexer and parser for PO (Portable Object) files.
// It includes token types for keywords like `msgid`, `msgstr`, `msgctxt`,
// and `msgid_plural`, as well as for comments, strings, and illegal or unexpected characters.
//
// Key Features:
// - Defines token types for PO file constructs.
// - Provides a function `LookupIdent` to identify and classify tokens based on their content.
// - Supports plural forms of message strings (`msgstr[...]`).
//
// Example Usage:
//
//	tok := token.Token{Type: token.MSGID, Literal: "Hello", Pos: 0}
//	fmt.Println(tok.Type) // Output: MSGID
//
// For more details, refer to the individual constants, types, and functions.
package token

import "regexp"

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	COMMENT = "COMMENT"
	STRING  = "STRING"
	MSGID   = "MSGID"
	MSGSTR  = "MSGSTR"
	MSGCTXT = "MSGCTXT"

	PluralMsgid  = "PluralMsgid"
	PluralMsgstr = "PluralMsgstr"
)

var pluralRegex = regexp.MustCompile(`msgstr\[*\d*\]`)

type Type string

var keywords = map[string]Type{
	"msgid":        MSGID,
	"msgstr":       MSGSTR,
	"msgctxt":      MSGCTXT,
	"msgid_plural": PluralMsgid,
	"#":            COMMENT,
}

type Token struct {
	Type    Type
	Literal string
	Pos     int
}

func LookupIdent(ident string) Type {
	if pluralRegex.MatchString(ident) {
		return PluralMsgstr
	}

	if token, ok := keywords[ident]; ok {
		return token
	}

	return ILLEGAL
}
