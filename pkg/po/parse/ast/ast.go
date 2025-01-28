package ast

import "fmt"

type Node interface {
	Pos() int
	Literal() string
	String() string
}

type (
	// Comment types.

	FlagComment struct {
		pos     int
		literal string
		Flag    string
	}
	LocationComment struct {
		pos     int
		literal string
		File    string
		Line    int
	}
	GeneralComment struct {
		pos     int
		literal string
		Text    string
	}

	// Identifiers.

	Msgid struct {
		pos     int
		literal string
		ID      string
	}
	Msgstr struct {
		pos     int
		literal string
		Str     string
	}
	Msgctxt struct {
		pos     int
		literal string
		Context string
	}

	// Plurals.

	MsgstrPlural struct {
		pos      int
		literal  string
		PluralID int
		Str      string
	}
	MsgidPlural struct {
		pos     int
		literal string
		Plural  string
	}

	File struct {
		pos     int
		literal string
		Name    string
		Nodes   []Node
	}
)

func (n File) Pos() int        { return n.pos }
func (n File) Literal() string { return n.literal }
func (n File) String() string  { return fmt.Sprintf("%#v", n) }

func (n FlagComment) Pos() int     { return n.pos }
func (n LocationComment) Pos() int { return n.pos }
func (n GeneralComment) Pos() int  { return n.pos }
func (n Msgid) Pos() int           { return n.pos }
func (n Msgstr) Pos() int          { return n.pos }
func (n Msgctxt) Pos() int         { return n.pos }
func (n MsgstrPlural) Pos() int    { return n.pos }
func (n MsgidPlural) Pos() int     { return n.pos }

func (n FlagComment) Literal() string     { return n.literal }
func (n LocationComment) Literal() string { return n.literal }
func (n GeneralComment) Literal() string  { return n.literal }
func (n Msgid) Literal() string           { return n.literal }
func (n Msgstr) Literal() string          { return n.literal }
func (n Msgctxt) Literal() string         { return n.literal }
func (n MsgstrPlural) Literal() string    { return n.literal }
func (n MsgidPlural) Literal() string     { return n.literal }

func (n FlagComment) String() string     { return FormatNode(n) }
func (n LocationComment) String() string { return FormatNode(n) }
func (n GeneralComment) String() string  { return FormatNode(n) }
func (n Msgid) String() string           { return FormatNode(n) }
func (n Msgstr) String() string          { return FormatNode(n) }
func (n Msgctxt) String() string         { return FormatNode(n) }
func (n MsgidPlural) String() string     { return FormatNode(n) }
func (n MsgstrPlural) String() string    { return FormatNode(n) }
