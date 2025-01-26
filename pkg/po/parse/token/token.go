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
