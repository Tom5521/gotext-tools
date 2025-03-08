// Code generated from ./internal/antlr-po/Po.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parse // Po
import "github.com/antlr4-go/antlr/v4"

// BasePoListener is a complete listener for a parse tree produced by PoParser.
type BasePoListener struct{}

var _ PoListener = &BasePoListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BasePoListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BasePoListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BasePoListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BasePoListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterStart is called when production start is entered.
func (s *BasePoListener) EnterStart(ctx *StartContext) {}

// ExitStart is called when production start is exited.
func (s *BasePoListener) ExitStart(ctx *StartContext) {}

// EnterEntry is called when production entry is entered.
func (s *BasePoListener) EnterEntry(ctx *EntryContext) {}

// ExitEntry is called when production entry is exited.
func (s *BasePoListener) ExitEntry(ctx *EntryContext) {}

// EnterMsgctxt is called when production msgctxt is entered.
func (s *BasePoListener) EnterMsgctxt(ctx *MsgctxtContext) {}

// ExitMsgctxt is called when production msgctxt is exited.
func (s *BasePoListener) ExitMsgctxt(ctx *MsgctxtContext) {}

// EnterMsgid is called when production msgid is entered.
func (s *BasePoListener) EnterMsgid(ctx *MsgidContext) {}

// ExitMsgid is called when production msgid is exited.
func (s *BasePoListener) ExitMsgid(ctx *MsgidContext) {}

// EnterMsgstr is called when production msgstr is entered.
func (s *BasePoListener) EnterMsgstr(ctx *MsgstrContext) {}

// ExitMsgstr is called when production msgstr is exited.
func (s *BasePoListener) ExitMsgstr(ctx *MsgstrContext) {}

// EnterPlural_msgid is called when production plural_msgid is entered.
func (s *BasePoListener) EnterPlural_msgid(ctx *Plural_msgidContext) {}

// ExitPlural_msgid is called when production plural_msgid is exited.
func (s *BasePoListener) ExitPlural_msgid(ctx *Plural_msgidContext) {}

// EnterPlural_msgstr is called when production plural_msgstr is entered.
func (s *BasePoListener) EnterPlural_msgstr(ctx *Plural_msgstrContext) {}

// ExitPlural_msgstr is called when production plural_msgstr is exited.
func (s *BasePoListener) ExitPlural_msgstr(ctx *Plural_msgstrContext) {}

// EnterString is called when production string is entered.
func (s *BasePoListener) EnterString(ctx *StringContext) {}

// ExitString is called when production string is exited.
func (s *BasePoListener) ExitString(ctx *StringContext) {}

// EnterComment is called when production comment is entered.
func (s *BasePoListener) EnterComment(ctx *CommentContext) {}

// ExitComment is called when production comment is exited.
func (s *BasePoListener) ExitComment(ctx *CommentContext) {}
