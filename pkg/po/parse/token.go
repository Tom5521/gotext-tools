package parse

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

type Type string

var keywords = map[string]Type{
	"msgid":        MSGID,
	"msgstr":       MSGSTR,
	"msgctxt":      MSGCTXT,
	"msgid_plural": PluralMsgid,
	"msgstr[%d]":   PluralMsgstr,
	"#":            COMMENT,
}

type Token struct {
	Type    Type
	Literal string
}

func LookupIdent(ident string) Type {
	if token, ok := keywords[ident]; ok {
		return token
	}

	return Type(ident)
}
