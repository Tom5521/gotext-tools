package parse

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/Tom5521/xgotext/pkg/po"
	"github.com/antlr4-go/antlr/v4"
)

var _ PoListener = (*Listener)(nil)

type Listener struct {
	*BasePoListener

	entry   po.Entry
	Entries po.Entries

	Errors []error
}

func (l *Listener) getStrings(strs []antlr.TerminalNode, storange *string) {
	var b strings.Builder

	for i, str := range strs {
		s, err := strconv.Unquote(str.GetText())
		if err != nil {
			l.Errors = append(l.Errors, err)
		}
		b.WriteString(s)
		if i != len(strs)-1 {
			b.WriteByte('\n')
		}
	}

	if storange != nil {
		*storange = b.String()
	}
}

func (l *Listener) ExitMsgctxt(ctx *MsgctxtContext) {
	l.getStrings(ctx.String_().AllSTRING(), &l.entry.Context)
}

func (l *Listener) EnterMsgid(ctx *MsgidContext) {
	l.getStrings(ctx.String_().AllSTRING(), &l.entry.ID)
}

func (l *Listener) EnterMsgstr(ctx *MsgstrContext) {
	l.getStrings(ctx.String_().AllSTRING(), &l.entry.Str)
}

func (l *Listener) EnterPlural_msgid(ctx *Plural_msgidContext) {
	l.getStrings(ctx.String_().AllSTRING(), &l.entry.Plural)
}

var pluralMsgstrRegex = regexp.MustCompile(`msgstr\[(\d*)\]`)

func (l *Listener) EnterPlural_msgstr(ctx *Plural_msgstrContext) {
	var plural po.PluralEntry
	l.getStrings(ctx.String_().AllSTRING(), &plural.Str)

	literal := ctx.GetText()
	literal = pluralMsgstrRegex.FindStringSubmatch(literal)[1]
	uintv, err := strconv.Atoi(literal)
	if err != nil {
		l.Errors = append(l.Errors, err)
		return
	}

	plural.ID = uintv
	l.entry.Plurals = append(l.entry.Plurals, plural)
}

func (l *Listener) ExitEntry(_ *EntryContext) {
	l.Entries = append(l.Entries, l.entry)
	l.entry = po.Entry{}
}

var (
	locationRegex  = regexp.MustCompile(`#:\s?(.*)`)
	generalRegex   = regexp.MustCompile(`#\s?(.*)`)
	extractedRegex = regexp.MustCompile(`#\.\s?(.*)`)
	flagRegex      = regexp.MustCompile(`#,\s?(.*)`)
	previousRegex  = regexp.MustCompile(`#\|\s?(.*)`)
)

func (l *Listener) EnterComment(ctx *CommentContext) {
	literal := ctx.GetText()

	switch {
	case locationRegex.MatchString(literal):
		matches := locationRegex.FindStringSubmatch(literal)
		parts := strings.SplitN(matches[1], ":", 2)
		line := -1
		var err error
		if parts[1] != "" {
			line, err = strconv.Atoi(parts[1])
			if err != nil {
				l.Errors = append(l.Errors, err)
				return
			}
		}

		loc := po.Location{
			Line: line,
			File: parts[0],
		}
		l.entry.Locations = append(l.entry.Locations, loc)
	case extractedRegex.MatchString(literal):
		l.entry.ExtractedComments = append(
			l.entry.ExtractedComments,
			extractedRegex.FindStringSubmatch(literal)[1],
		)
	case flagRegex.MatchString(literal):
		l.entry.Flags = append(l.entry.Flags, flagRegex.FindStringSubmatch(literal)[1])
	case previousRegex.MatchString(literal):
		l.entry.Previous = append(l.entry.Previous, previousRegex.FindStringSubmatch(literal)[1])
	default:
		l.entry.Comments = append(l.entry.Comments, generalRegex.FindStringSubmatch(literal)[1])
	}
}

func (l *Listener) VisitErrorNode(node antlr.ErrorNode) {
	l.Errors = append(l.Errors, errors.New(node.GetText()))
}
