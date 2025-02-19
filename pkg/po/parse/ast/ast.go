package ast

type Node interface {
	Pos() int
}

type (
	Entry struct {
		pos               int
		Flags             []*FlagComment
		ExtractedComments []*ExtractedComment
		LocationComments  []*LocationComment
		GeneralComments   []*GeneralComment
		PreviousComments  []*PreviousComment
		Msgid             *Msgid
		Msgstr            *Msgstr
		Msgctxt           *Msgctxt
		Plural            *MsgidPlural
		Plurals           []*MsgstrPlural
	}

	// Comment types.
	PluralEntry struct {
		ID  int
		Str string
	}

	FlagComment struct {
		pos  int
		Flag string
	}
	LocationComment struct {
		pos  int
		File string
		Line int
	}
	GeneralComment struct {
		pos  int
		Text string
	}
	ExtractedComment struct {
		pos  int
		Text string
	}
	PreviousComment struct {
		pos  int
		Text string
	}

	// Identifiers.

	Msgid struct {
		pos int
		ID  string
	}
	Msgstr struct {
		pos int
		Str string
	}
	Msgctxt struct {
		pos     int
		Context string
	}

	// Plurals.

	MsgstrPlural struct {
		pos      int
		PluralID int
		Str      string
	}
	MsgidPlural struct {
		pos    int
		Plural string
	}

	AST struct {
		pos     int
		Name    string
		Content []byte
		Nodes   []Node
	}
)

func (e PreviousComment) Pos() int  { return e.pos }
func (e ExtractedComment) Pos() int { return e.pos }
func (e Entry) Pos() int            { return e.pos }
func (n AST) Pos() int              { return n.pos }
func (n FlagComment) Pos() int      { return n.pos }
func (n LocationComment) Pos() int  { return n.pos }
func (n GeneralComment) Pos() int   { return n.pos }
func (n Msgid) Pos() int            { return n.pos }
func (n Msgstr) Pos() int           { return n.pos }
func (n Msgctxt) Pos() int          { return n.pos }
func (n MsgstrPlural) Pos() int     { return n.pos }
func (n MsgidPlural) Pos() int      { return n.pos }
