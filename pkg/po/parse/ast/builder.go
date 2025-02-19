package ast

import (
	"errors"
	"fmt"
	"reflect"
	"slices"

	"github.com/Tom5521/xgotext/internal/util"
)

func tfor[T any]() reflect.Type {
	return reflect.TypeFor[T]()
}

func types(a ...any) []reflect.Type {
	var t []reflect.Type
	for _, b := range a {
		t = append(t, reflect.TypeOf(b))
	}
	return t
}

type ASTBuilder struct {
	curEntry Entry

	input []Node

	name    string
	content []byte

	entries []Node
	warns   []string
	errs    []error

	toSkip []reflect.Type
}

func NewASTBuilder(name string, content []byte, nodes []Node) *ASTBuilder {
	n := &ASTBuilder{
		name:    name,
		content: content,
		input:   nodes,
	}
	n.reset()

	return n
}

func (n *ASTBuilder) finishEntry(cur Node) {
	err := validateEntry(n.curEntry)
	if err != nil {
		n.appendErr(err)
		return
	}

	foundID := n.curEntry.Msgid != nil
	foundStr := n.curEntry.Msgstr != nil || len(n.curEntry.Plurals) > 0

	if !foundID {
		n.appendErrf(
			"msgid not found at %s:%d",
			n.name,
			util.FindLine(n.content, cur.Pos()),
		)

		return
	}

	if !foundStr {
		n.appendWarn(
			"msgstr not found at %s:%d",
			n.name,
			util.FindLine(n.content, cur.Pos()),
		)
	}

	if !foundStr || !foundID {
		return
	}

	n.entries = append(n.entries, n.curEntry)
	n.resetState()
}

func (n *ASTBuilder) resetState() {
	n.curEntry = Entry{
		pos: -1,
	}
}

func (n *ASTBuilder) reset() {
	n.toSkip = types(
		FlagComment{},
		GeneralComment{},
		ExtractedComment{},
		LocationComment{},
		PreviousComment{},
	)
	n.entries = nil
	n.warns = nil
	n.errs = nil
	n.resetState()
}

func (n *ASTBuilder) AST() *AST {
	return &AST{
		pos:     0,
		Content: n.content,
		Name:    n.name,
		Nodes:   n.entries,
	}
}

func (n *ASTBuilder) Entries() []Entry {
	var entries []Entry
	for _, node := range n.entries {
		entries = append(entries, node.(Entry))
	}

	return entries
}

func (n *ASTBuilder) Errors() []error {
	return n.errs
}

func (n *ASTBuilder) Warnings() []string {
	return n.warns
}

func (n *ASTBuilder) genParseMap() map[reflect.Type]func(Node, int) {
	return map[reflect.Type]func(Node, int){
		// Comments.
		tfor[GeneralComment]():   n.handleComment,
		tfor[FlagComment]():      n.handleComment,
		tfor[LocationComment]():  n.handleComment,
		tfor[ExtractedComment](): n.handleComment,
		tfor[PreviousComment]():  n.handleComment,

		tfor[Msgctxt]():      n.handleMsgctxt,
		tfor[Msgid]():        n.handleMsgid,
		tfor[Msgstr]():       n.handleMsgstr,
		tfor[MsgstrPlural](): n.handleMsgstrPlural,
		tfor[MsgidPlural]():  n.handleMsgidPlural,
	}
}

func (n *ASTBuilder) handleComment(node Node, i int) {
	switch node := node.(type) {
	case LocationComment:
		n.curEntry.LocationComments = append(n.curEntry.LocationComments, &node)
	case GeneralComment:
		n.curEntry.GeneralComments = append(n.curEntry.GeneralComments, &node)
	case ExtractedComment:
		n.curEntry.ExtractedComments = append(n.curEntry.ExtractedComments, &node)
	case FlagComment:
		n.curEntry.Flags = append(n.curEntry.Flags, &node)
	case PreviousComment:
		n.curEntry.PreviousComments = append(n.curEntry.PreviousComments, &node)
	}
}

func (n *ASTBuilder) handleMsgctxt(node Node, i int) {
	msgctxt := node.(Msgctxt)

	if n.curEntry.pos == -1 {
		n.curEntry.pos = node.Pos()
	}

	if n.curEntry.Msgctxt != nil {
		n.appendWarn(
			"duplicated msgctxt at %s:%d",
			n.name,
			util.FindLine(n.content, node.Pos()),
		)
	}

	n.curEntry.Msgctxt = &msgctxt

	if !n.typeIsComing(i+1, n.toSkip, tfor[Msgid]()) {
		n.finishEntry(node)
	}
}

func (n *ASTBuilder) handleMsgid(node Node, i int) {
	msgid := node.(Msgid)

	if n.curEntry.pos == -1 {
		n.curEntry.pos = node.Pos()
	}

	if n.curEntry.Msgid != nil {
		n.appendWarn(
			"duplicated msgid at %s:%d",
			n.name,
			util.FindLine(n.content, node.Pos()),
		)
	}

	n.curEntry.Msgid = &msgid

	if !n.typeIsComing(i+1, n.toSkip, tfor[Msgstr](), tfor[MsgidPlural](), tfor[MsgstrPlural]()) {
		n.finishEntry(node)
	}
}

func (n *ASTBuilder) handleMsgstr(node Node, i int) {
	msgstr := node.(Msgstr)

	if n.curEntry.Plural != nil {
		return
	}

	n.curEntry.Msgstr = &msgstr
	n.finishEntry(node)
}

func (n *ASTBuilder) handleMsgidPlural(node Node, i int) {
	msgidPlural := node.(MsgidPlural)

	if n.curEntry.Plural != nil {
		n.appendWarn(
			"duplicated msgid_plural at %s:%d",
			n.name,
			util.FindLine(n.content, node.Pos()),
		)
	}

	n.curEntry.Plural = &msgidPlural

	if !n.typeIsComing(i+1, n.toSkip, tfor[MsgstrPlural]()) {
		n.finishEntry(node)
	}
}

func (n *ASTBuilder) handleMsgstrPlural(node Node, i int) {
	msgstrPlural := node.(MsgstrPlural)

	n.curEntry.Plurals = append(n.curEntry.Plurals, &msgstrPlural)

	if !n.typeIsComing(i+1, n.toSkip, tfor[MsgstrPlural]()) {
		n.finishEntry(node)
	}
}

func (n *ASTBuilder) appendWarn(format string, a ...any) {
	n.warns = append(n.warns, fmt.Sprintf(format, a...))
}

func (n *ASTBuilder) appendErr(err error) {
	n.errs = append(n.errs, err)
}

func (n *ASTBuilder) appendErrf(format string, a ...any) {
	n.appendErr(fmt.Errorf(format, a...))
}

func (n *ASTBuilder) comingType(offset int, ignore []reflect.Type) reflect.Type {
	for _, node := range n.input[offset:] {
		t := reflect.TypeOf(node)
		if slices.Contains(ignore, t) {
			continue
		}

		return t
	}

	return nil
}

func (n *ASTBuilder) typeIsComing(offset int, ignore []reflect.Type, wanted ...reflect.Type) bool {
	t := n.comingType(offset, ignore)
	return slices.Contains(wanted, t)
}

func (n *ASTBuilder) Build() {
	n.reset()

	parseMap := n.genParseMap()

	for i, node := range n.input {
		parseMap[reflect.TypeOf(node)](node, i)
	}
}

func validateEntry(e Entry) error {
	if e.Plural == nil && len(e.Plurals) > 0 {
		return errors.New("plural translations provided but no plural form has been specified")
	}

	if e.Msgstr == nil && len(e.Plurals) == 0 {
		return errors.New("no msgstr was specified")
	}

	if e.Plurals != nil && len(e.Plurals) == 0 {
		return errors.New("plural form specified but no plural translations provided")
	}

	return nil
}
