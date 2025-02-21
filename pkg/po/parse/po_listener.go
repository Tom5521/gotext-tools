// Code generated from ./pkg/po/parse/Po.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parse // Po
import "github.com/antlr4-go/antlr/v4"

// PoListener is a complete listener for a parse tree produced by PoParser.
type PoListener interface {
	antlr.ParseTreeListener

	// EnterStart is called when entering the start production.
	EnterStart(c *StartContext)

	// EnterEntry is called when entering the entry production.
	EnterEntry(c *EntryContext)

	// EnterMsgctxt is called when entering the msgctxt production.
	EnterMsgctxt(c *MsgctxtContext)

	// EnterMsgid is called when entering the msgid production.
	EnterMsgid(c *MsgidContext)

	// EnterMsgstr is called when entering the msgstr production.
	EnterMsgstr(c *MsgstrContext)

	// EnterPlural_msgid is called when entering the plural_msgid production.
	EnterPlural_msgid(c *Plural_msgidContext)

	// EnterPlural_msgstr is called when entering the plural_msgstr production.
	EnterPlural_msgstr(c *Plural_msgstrContext)

	// EnterString is called when entering the string production.
	EnterString(c *StringContext)

	// EnterComment is called when entering the comment production.
	EnterComment(c *CommentContext)

	// ExitStart is called when exiting the start production.
	ExitStart(c *StartContext)

	// ExitEntry is called when exiting the entry production.
	ExitEntry(c *EntryContext)

	// ExitMsgctxt is called when exiting the msgctxt production.
	ExitMsgctxt(c *MsgctxtContext)

	// ExitMsgid is called when exiting the msgid production.
	ExitMsgid(c *MsgidContext)

	// ExitMsgstr is called when exiting the msgstr production.
	ExitMsgstr(c *MsgstrContext)

	// ExitPlural_msgid is called when exiting the plural_msgid production.
	ExitPlural_msgid(c *Plural_msgidContext)

	// ExitPlural_msgstr is called when exiting the plural_msgstr production.
	ExitPlural_msgstr(c *Plural_msgstrContext)

	// ExitString is called when exiting the string production.
	ExitString(c *StringContext)

	// ExitComment is called when exiting the comment production.
	ExitComment(c *CommentContext)
}
