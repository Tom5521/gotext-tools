// Code generated from ./pkg/po/parse/Po.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parse // Po
import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type PoParser struct {
	*antlr.BaseParser
}

var PoParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func poParserInit() {
	staticData := &PoParserStaticData
	staticData.LiteralNames = []string{
		"", "", "", "", "'\\n'", "'msgctxt'", "'msgid'", "'msgstr'", "'msgid_plural'",
	}
	staticData.SymbolicNames = []string{
		"", "WS", "INT", "STRING", "NL", "MSGCTXT", "MSGID", "MSGSTR", "PLURAL_MSGID",
		"PLURAL_MSGSTR", "COMMENT", "FLAG_COMMENT", "EXTRACTED_COMMENT", "REFERENCE_COMMENT",
		"PREVIOUS_COMMENT",
	}
	staticData.RuleNames = []string{
		"start", "entry", "msgctxt", "msgid", "msgstr", "plural_msgid", "plural_msgstr",
		"string", "comment",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 14, 72, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 1, 0, 5, 0, 20, 8, 0,
		10, 0, 12, 0, 23, 9, 0, 1, 1, 5, 1, 26, 8, 1, 10, 1, 12, 1, 29, 9, 1, 1,
		1, 3, 1, 32, 8, 1, 1, 1, 1, 1, 1, 1, 1, 1, 4, 1, 38, 8, 1, 11, 1, 12, 1,
		39, 3, 1, 42, 8, 1, 1, 1, 3, 1, 45, 8, 1, 1, 2, 1, 2, 1, 2, 1, 3, 1, 3,
		1, 3, 1, 4, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 1, 7, 1, 7,
		3, 7, 64, 8, 7, 4, 7, 66, 8, 7, 11, 7, 12, 7, 67, 1, 8, 1, 8, 1, 8, 0,
		0, 9, 0, 2, 4, 6, 8, 10, 12, 14, 16, 0, 1, 1, 0, 10, 14, 70, 0, 21, 1,
		0, 0, 0, 2, 27, 1, 0, 0, 0, 4, 46, 1, 0, 0, 0, 6, 49, 1, 0, 0, 0, 8, 52,
		1, 0, 0, 0, 10, 55, 1, 0, 0, 0, 12, 58, 1, 0, 0, 0, 14, 65, 1, 0, 0, 0,
		16, 69, 1, 0, 0, 0, 18, 20, 3, 2, 1, 0, 19, 18, 1, 0, 0, 0, 20, 23, 1,
		0, 0, 0, 21, 19, 1, 0, 0, 0, 21, 22, 1, 0, 0, 0, 22, 1, 1, 0, 0, 0, 23,
		21, 1, 0, 0, 0, 24, 26, 3, 16, 8, 0, 25, 24, 1, 0, 0, 0, 26, 29, 1, 0,
		0, 0, 27, 25, 1, 0, 0, 0, 27, 28, 1, 0, 0, 0, 28, 31, 1, 0, 0, 0, 29, 27,
		1, 0, 0, 0, 30, 32, 3, 4, 2, 0, 31, 30, 1, 0, 0, 0, 31, 32, 1, 0, 0, 0,
		32, 33, 1, 0, 0, 0, 33, 41, 3, 6, 3, 0, 34, 42, 3, 8, 4, 0, 35, 37, 3,
		10, 5, 0, 36, 38, 3, 12, 6, 0, 37, 36, 1, 0, 0, 0, 38, 39, 1, 0, 0, 0,
		39, 37, 1, 0, 0, 0, 39, 40, 1, 0, 0, 0, 40, 42, 1, 0, 0, 0, 41, 34, 1,
		0, 0, 0, 41, 35, 1, 0, 0, 0, 42, 44, 1, 0, 0, 0, 43, 45, 5, 4, 0, 0, 44,
		43, 1, 0, 0, 0, 44, 45, 1, 0, 0, 0, 45, 3, 1, 0, 0, 0, 46, 47, 5, 5, 0,
		0, 47, 48, 3, 14, 7, 0, 48, 5, 1, 0, 0, 0, 49, 50, 5, 6, 0, 0, 50, 51,
		3, 14, 7, 0, 51, 7, 1, 0, 0, 0, 52, 53, 5, 7, 0, 0, 53, 54, 3, 14, 7, 0,
		54, 9, 1, 0, 0, 0, 55, 56, 5, 8, 0, 0, 56, 57, 3, 14, 7, 0, 57, 11, 1,
		0, 0, 0, 58, 59, 5, 9, 0, 0, 59, 60, 3, 14, 7, 0, 60, 13, 1, 0, 0, 0, 61,
		63, 5, 3, 0, 0, 62, 64, 5, 4, 0, 0, 63, 62, 1, 0, 0, 0, 63, 64, 1, 0, 0,
		0, 64, 66, 1, 0, 0, 0, 65, 61, 1, 0, 0, 0, 66, 67, 1, 0, 0, 0, 67, 65,
		1, 0, 0, 0, 67, 68, 1, 0, 0, 0, 68, 15, 1, 0, 0, 0, 69, 70, 7, 0, 0, 0,
		70, 17, 1, 0, 0, 0, 8, 21, 27, 31, 39, 41, 44, 63, 67,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// PoParserInit initializes any static state used to implement PoParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewPoParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func PoParserInit() {
	staticData := &PoParserStaticData
	staticData.once.Do(poParserInit)
}

// NewPoParser produces a new parser instance for the optional input antlr.TokenStream.
func NewPoParser(input antlr.TokenStream) *PoParser {
	PoParserInit()
	this := new(PoParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &PoParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "Po.g4"

	return this
}

// PoParser tokens.
const (
	PoParserEOF               = antlr.TokenEOF
	PoParserWS                = 1
	PoParserINT               = 2
	PoParserSTRING            = 3
	PoParserNL                = 4
	PoParserMSGCTXT           = 5
	PoParserMSGID             = 6
	PoParserMSGSTR            = 7
	PoParserPLURAL_MSGID      = 8
	PoParserPLURAL_MSGSTR     = 9
	PoParserCOMMENT           = 10
	PoParserFLAG_COMMENT      = 11
	PoParserEXTRACTED_COMMENT = 12
	PoParserREFERENCE_COMMENT = 13
	PoParserPREVIOUS_COMMENT  = 14
)

// PoParser rules.
const (
	PoParserRULE_start         = 0
	PoParserRULE_entry         = 1
	PoParserRULE_msgctxt       = 2
	PoParserRULE_msgid         = 3
	PoParserRULE_msgstr        = 4
	PoParserRULE_plural_msgid  = 5
	PoParserRULE_plural_msgstr = 6
	PoParserRULE_string        = 7
	PoParserRULE_comment       = 8
)

// IStartContext is an interface to support dynamic dispatch.
type IStartContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllEntry() []IEntryContext
	Entry(i int) IEntryContext

	// IsStartContext differentiates from other interfaces.
	IsStartContext()
}

type StartContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStartContext() *StartContext {
	var p = new(StartContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PoParserRULE_start
	return p
}

func InitEmptyStartContext(p *StartContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PoParserRULE_start
}

func (*StartContext) IsStartContext() {}

func NewStartContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StartContext {
	var p = new(StartContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PoParserRULE_start

	return p
}

func (s *StartContext) GetParser() antlr.Parser { return s.parser }

func (s *StartContext) AllEntry() []IEntryContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IEntryContext); ok {
			len++
		}
	}

	tst := make([]IEntryContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IEntryContext); ok {
			tst[i] = t.(IEntryContext)
			i++
		}
	}

	return tst
}

func (s *StartContext) Entry(i int) IEntryContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEntryContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEntryContext)
}

func (s *StartContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StartContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StartContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PoListener); ok {
		listenerT.EnterStart(s)
	}
}

func (s *StartContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PoListener); ok {
		listenerT.ExitStart(s)
	}
}

func (p *PoParser) Start_() (localctx IStartContext) {
	localctx = NewStartContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, PoParserRULE_start)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(21)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&31840) != 0 {
		{
			p.SetState(18)
			p.Entry()
		}

		p.SetState(23)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IEntryContext is an interface to support dynamic dispatch.
type IEntryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Msgid() IMsgidContext
	Msgstr() IMsgstrContext
	AllComment() []ICommentContext
	Comment(i int) ICommentContext
	Msgctxt() IMsgctxtContext
	NL() antlr.TerminalNode
	Plural_msgid() IPlural_msgidContext
	AllPlural_msgstr() []IPlural_msgstrContext
	Plural_msgstr(i int) IPlural_msgstrContext

	// IsEntryContext differentiates from other interfaces.
	IsEntryContext()
}

type EntryContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEntryContext() *EntryContext {
	var p = new(EntryContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PoParserRULE_entry
	return p
}

func InitEmptyEntryContext(p *EntryContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PoParserRULE_entry
}

func (*EntryContext) IsEntryContext() {}

func NewEntryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EntryContext {
	var p = new(EntryContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PoParserRULE_entry

	return p
}

func (s *EntryContext) GetParser() antlr.Parser { return s.parser }

func (s *EntryContext) Msgid() IMsgidContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMsgidContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMsgidContext)
}

func (s *EntryContext) Msgstr() IMsgstrContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMsgstrContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMsgstrContext)
}

func (s *EntryContext) AllComment() []ICommentContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ICommentContext); ok {
			len++
		}
	}

	tst := make([]ICommentContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ICommentContext); ok {
			tst[i] = t.(ICommentContext)
			i++
		}
	}

	return tst
}

func (s *EntryContext) Comment(i int) ICommentContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICommentContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICommentContext)
}

func (s *EntryContext) Msgctxt() IMsgctxtContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMsgctxtContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMsgctxtContext)
}

func (s *EntryContext) NL() antlr.TerminalNode {
	return s.GetToken(PoParserNL, 0)
}

func (s *EntryContext) Plural_msgid() IPlural_msgidContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPlural_msgidContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPlural_msgidContext)
}

func (s *EntryContext) AllPlural_msgstr() []IPlural_msgstrContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IPlural_msgstrContext); ok {
			len++
		}
	}

	tst := make([]IPlural_msgstrContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IPlural_msgstrContext); ok {
			tst[i] = t.(IPlural_msgstrContext)
			i++
		}
	}

	return tst
}

func (s *EntryContext) Plural_msgstr(i int) IPlural_msgstrContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPlural_msgstrContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPlural_msgstrContext)
}

func (s *EntryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EntryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EntryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PoListener); ok {
		listenerT.EnterEntry(s)
	}
}

func (s *EntryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PoListener); ok {
		listenerT.ExitEntry(s)
	}
}

func (p *PoParser) Entry() (localctx IEntryContext) {
	localctx = NewEntryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, PoParserRULE_entry)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(27)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&31744) != 0 {
		{
			p.SetState(24)
			p.Comment()
		}

		p.SetState(29)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(31)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == PoParserMSGCTXT {
		{
			p.SetState(30)
			p.Msgctxt()
		}

	}
	{
		p.SetState(33)
		p.Msgid()
	}
	p.SetState(41)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case PoParserMSGSTR:
		{
			p.SetState(34)
			p.Msgstr()
		}

	case PoParserPLURAL_MSGID:
		{
			p.SetState(35)
			p.Plural_msgid()
		}
		p.SetState(37)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == PoParserPLURAL_MSGSTR {
			{
				p.SetState(36)
				p.Plural_msgstr()
			}

			p.SetState(39)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}
	p.SetState(44)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == PoParserNL {
		{
			p.SetState(43)
			p.Match(PoParserNL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IMsgctxtContext is an interface to support dynamic dispatch.
type IMsgctxtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	MSGCTXT() antlr.TerminalNode
	String_() IStringContext

	// IsMsgctxtContext differentiates from other interfaces.
	IsMsgctxtContext()
}

type MsgctxtContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMsgctxtContext() *MsgctxtContext {
	var p = new(MsgctxtContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PoParserRULE_msgctxt
	return p
}

func InitEmptyMsgctxtContext(p *MsgctxtContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PoParserRULE_msgctxt
}

func (*MsgctxtContext) IsMsgctxtContext() {}

func NewMsgctxtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MsgctxtContext {
	var p = new(MsgctxtContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PoParserRULE_msgctxt

	return p
}

func (s *MsgctxtContext) GetParser() antlr.Parser { return s.parser }

func (s *MsgctxtContext) MSGCTXT() antlr.TerminalNode {
	return s.GetToken(PoParserMSGCTXT, 0)
}

func (s *MsgctxtContext) String_() IStringContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStringContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStringContext)
}

func (s *MsgctxtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MsgctxtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MsgctxtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PoListener); ok {
		listenerT.EnterMsgctxt(s)
	}
}

func (s *MsgctxtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PoListener); ok {
		listenerT.ExitMsgctxt(s)
	}
}

func (p *PoParser) Msgctxt() (localctx IMsgctxtContext) {
	localctx = NewMsgctxtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, PoParserRULE_msgctxt)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(46)
		p.Match(PoParserMSGCTXT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(47)
		p.String_()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IMsgidContext is an interface to support dynamic dispatch.
type IMsgidContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	MSGID() antlr.TerminalNode
	String_() IStringContext

	// IsMsgidContext differentiates from other interfaces.
	IsMsgidContext()
}

type MsgidContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMsgidContext() *MsgidContext {
	var p = new(MsgidContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PoParserRULE_msgid
	return p
}

func InitEmptyMsgidContext(p *MsgidContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PoParserRULE_msgid
}

func (*MsgidContext) IsMsgidContext() {}

func NewMsgidContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MsgidContext {
	var p = new(MsgidContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PoParserRULE_msgid

	return p
}

func (s *MsgidContext) GetParser() antlr.Parser { return s.parser }

func (s *MsgidContext) MSGID() antlr.TerminalNode {
	return s.GetToken(PoParserMSGID, 0)
}

func (s *MsgidContext) String_() IStringContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStringContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStringContext)
}

func (s *MsgidContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MsgidContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MsgidContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PoListener); ok {
		listenerT.EnterMsgid(s)
	}
}

func (s *MsgidContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PoListener); ok {
		listenerT.ExitMsgid(s)
	}
}

func (p *PoParser) Msgid() (localctx IMsgidContext) {
	localctx = NewMsgidContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, PoParserRULE_msgid)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(49)
		p.Match(PoParserMSGID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(50)
		p.String_()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IMsgstrContext is an interface to support dynamic dispatch.
type IMsgstrContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	MSGSTR() antlr.TerminalNode
	String_() IStringContext

	// IsMsgstrContext differentiates from other interfaces.
	IsMsgstrContext()
}

type MsgstrContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMsgstrContext() *MsgstrContext {
	var p = new(MsgstrContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PoParserRULE_msgstr
	return p
}

func InitEmptyMsgstrContext(p *MsgstrContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PoParserRULE_msgstr
}

func (*MsgstrContext) IsMsgstrContext() {}

func NewMsgstrContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MsgstrContext {
	var p = new(MsgstrContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PoParserRULE_msgstr

	return p
}

func (s *MsgstrContext) GetParser() antlr.Parser { return s.parser }

func (s *MsgstrContext) MSGSTR() antlr.TerminalNode {
	return s.GetToken(PoParserMSGSTR, 0)
}

func (s *MsgstrContext) String_() IStringContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStringContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStringContext)
}

func (s *MsgstrContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MsgstrContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MsgstrContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PoListener); ok {
		listenerT.EnterMsgstr(s)
	}
}

func (s *MsgstrContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PoListener); ok {
		listenerT.ExitMsgstr(s)
	}
}

func (p *PoParser) Msgstr() (localctx IMsgstrContext) {
	localctx = NewMsgstrContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, PoParserRULE_msgstr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(52)
		p.Match(PoParserMSGSTR)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(53)
		p.String_()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPlural_msgidContext is an interface to support dynamic dispatch.
type IPlural_msgidContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PLURAL_MSGID() antlr.TerminalNode
	String_() IStringContext

	// IsPlural_msgidContext differentiates from other interfaces.
	IsPlural_msgidContext()
}

type Plural_msgidContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPlural_msgidContext() *Plural_msgidContext {
	var p = new(Plural_msgidContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PoParserRULE_plural_msgid
	return p
}

func InitEmptyPlural_msgidContext(p *Plural_msgidContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PoParserRULE_plural_msgid
}

func (*Plural_msgidContext) IsPlural_msgidContext() {}

func NewPlural_msgidContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Plural_msgidContext {
	var p = new(Plural_msgidContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PoParserRULE_plural_msgid

	return p
}

func (s *Plural_msgidContext) GetParser() antlr.Parser { return s.parser }

func (s *Plural_msgidContext) PLURAL_MSGID() antlr.TerminalNode {
	return s.GetToken(PoParserPLURAL_MSGID, 0)
}

func (s *Plural_msgidContext) String_() IStringContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStringContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStringContext)
}

func (s *Plural_msgidContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Plural_msgidContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Plural_msgidContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PoListener); ok {
		listenerT.EnterPlural_msgid(s)
	}
}

func (s *Plural_msgidContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PoListener); ok {
		listenerT.ExitPlural_msgid(s)
	}
}

func (p *PoParser) Plural_msgid() (localctx IPlural_msgidContext) {
	localctx = NewPlural_msgidContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, PoParserRULE_plural_msgid)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(55)
		p.Match(PoParserPLURAL_MSGID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(56)
		p.String_()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPlural_msgstrContext is an interface to support dynamic dispatch.
type IPlural_msgstrContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PLURAL_MSGSTR() antlr.TerminalNode
	String_() IStringContext

	// IsPlural_msgstrContext differentiates from other interfaces.
	IsPlural_msgstrContext()
}

type Plural_msgstrContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPlural_msgstrContext() *Plural_msgstrContext {
	var p = new(Plural_msgstrContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PoParserRULE_plural_msgstr
	return p
}

func InitEmptyPlural_msgstrContext(p *Plural_msgstrContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PoParserRULE_plural_msgstr
}

func (*Plural_msgstrContext) IsPlural_msgstrContext() {}

func NewPlural_msgstrContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Plural_msgstrContext {
	var p = new(Plural_msgstrContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PoParserRULE_plural_msgstr

	return p
}

func (s *Plural_msgstrContext) GetParser() antlr.Parser { return s.parser }

func (s *Plural_msgstrContext) PLURAL_MSGSTR() antlr.TerminalNode {
	return s.GetToken(PoParserPLURAL_MSGSTR, 0)
}

func (s *Plural_msgstrContext) String_() IStringContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStringContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStringContext)
}

func (s *Plural_msgstrContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Plural_msgstrContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Plural_msgstrContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PoListener); ok {
		listenerT.EnterPlural_msgstr(s)
	}
}

func (s *Plural_msgstrContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PoListener); ok {
		listenerT.ExitPlural_msgstr(s)
	}
}

func (p *PoParser) Plural_msgstr() (localctx IPlural_msgstrContext) {
	localctx = NewPlural_msgstrContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, PoParserRULE_plural_msgstr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(58)
		p.Match(PoParserPLURAL_MSGSTR)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(59)
		p.String_()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IStringContext is an interface to support dynamic dispatch.
type IStringContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllSTRING() []antlr.TerminalNode
	STRING(i int) antlr.TerminalNode
	AllNL() []antlr.TerminalNode
	NL(i int) antlr.TerminalNode

	// IsStringContext differentiates from other interfaces.
	IsStringContext()
}

type StringContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStringContext() *StringContext {
	var p = new(StringContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PoParserRULE_string
	return p
}

func InitEmptyStringContext(p *StringContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PoParserRULE_string
}

func (*StringContext) IsStringContext() {}

func NewStringContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StringContext {
	var p = new(StringContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PoParserRULE_string

	return p
}

func (s *StringContext) GetParser() antlr.Parser { return s.parser }

func (s *StringContext) AllSTRING() []antlr.TerminalNode {
	return s.GetTokens(PoParserSTRING)
}

func (s *StringContext) STRING(i int) antlr.TerminalNode {
	return s.GetToken(PoParserSTRING, i)
}

func (s *StringContext) AllNL() []antlr.TerminalNode {
	return s.GetTokens(PoParserNL)
}

func (s *StringContext) NL(i int) antlr.TerminalNode {
	return s.GetToken(PoParserNL, i)
}

func (s *StringContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StringContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StringContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PoListener); ok {
		listenerT.EnterString(s)
	}
}

func (s *StringContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PoListener); ok {
		listenerT.ExitString(s)
	}
}

func (p *PoParser) String_() (localctx IStringContext) {
	localctx = NewStringContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, PoParserRULE_string)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(65)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == PoParserSTRING {
		{
			p.SetState(61)
			p.Match(PoParserSTRING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(63)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 6, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(62)
				p.Match(PoParserNL)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		} else if p.HasError() { // JIM
			goto errorExit
		}

		p.SetState(67)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICommentContext is an interface to support dynamic dispatch.
type ICommentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	COMMENT() antlr.TerminalNode
	FLAG_COMMENT() antlr.TerminalNode
	EXTRACTED_COMMENT() antlr.TerminalNode
	REFERENCE_COMMENT() antlr.TerminalNode
	PREVIOUS_COMMENT() antlr.TerminalNode

	// IsCommentContext differentiates from other interfaces.
	IsCommentContext()
}

type CommentContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCommentContext() *CommentContext {
	var p = new(CommentContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PoParserRULE_comment
	return p
}

func InitEmptyCommentContext(p *CommentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PoParserRULE_comment
}

func (*CommentContext) IsCommentContext() {}

func NewCommentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CommentContext {
	var p = new(CommentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PoParserRULE_comment

	return p
}

func (s *CommentContext) GetParser() antlr.Parser { return s.parser }

func (s *CommentContext) COMMENT() antlr.TerminalNode {
	return s.GetToken(PoParserCOMMENT, 0)
}

func (s *CommentContext) FLAG_COMMENT() antlr.TerminalNode {
	return s.GetToken(PoParserFLAG_COMMENT, 0)
}

func (s *CommentContext) EXTRACTED_COMMENT() antlr.TerminalNode {
	return s.GetToken(PoParserEXTRACTED_COMMENT, 0)
}

func (s *CommentContext) REFERENCE_COMMENT() antlr.TerminalNode {
	return s.GetToken(PoParserREFERENCE_COMMENT, 0)
}

func (s *CommentContext) PREVIOUS_COMMENT() antlr.TerminalNode {
	return s.GetToken(PoParserPREVIOUS_COMMENT, 0)
}

func (s *CommentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CommentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CommentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PoListener); ok {
		listenerT.EnterComment(s)
	}
}

func (s *CommentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PoListener); ok {
		listenerT.ExitComment(s)
	}
}

func (p *PoParser) Comment() (localctx ICommentContext) {
	localctx = NewCommentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, PoParserRULE_comment)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(69)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&31744) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}
