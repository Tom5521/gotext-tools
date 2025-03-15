// Code generated from ./internal/antlr-po/Po.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parse

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"sync"
	"unicode"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter

type PoLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var PoLexerLexerStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	ChannelNames           []string
	ModeNames              []string
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func polexerLexerInit() {
	staticData := &PoLexerLexerStaticData
	staticData.ChannelNames = []string{
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
	}
	staticData.ModeNames = []string{
		"DEFAULT_MODE",
	}
	staticData.LiteralNames = []string{
		"", "", "", "", "'\\n'", "'msgctxt'", "'msgid'", "'msgstr'", "'msgid_plural'",
	}
	staticData.SymbolicNames = []string{
		"", "WS", "INT", "STRING", "NL", "MSGCTXT", "MSGID", "MSGSTR", "PLURAL_MSGID",
		"PLURAL_MSGSTR", "COMMENT",
	}
	staticData.RuleNames = []string{
		"WS", "INT", "STRING", "NL", "MSGCTXT", "MSGID", "MSGSTR", "PLURAL_MSGID",
		"PLURAL_MSGSTR", "COMMENT",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 10, 105, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 4, 1, 27, 8, 1, 11, 1, 12, 1, 28, 1, 2, 1, 2,
		1, 2, 1, 2, 5, 2, 35, 8, 2, 10, 2, 12, 2, 38, 9, 2, 1, 2, 1, 2, 1, 2, 1,
		2, 1, 2, 5, 2, 45, 8, 2, 10, 2, 12, 2, 48, 9, 2, 3, 2, 50, 8, 2, 1, 3,
		1, 3, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5,
		1, 5, 1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 7, 1, 7,
		1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 8,
		1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 9, 1, 9,
		5, 9, 101, 8, 9, 10, 9, 12, 9, 104, 9, 9, 0, 0, 10, 1, 1, 3, 2, 5, 3, 7,
		4, 9, 5, 11, 6, 13, 7, 15, 8, 17, 9, 19, 10, 1, 0, 5, 3, 0, 9, 10, 13,
		13, 32, 32, 1, 0, 48, 57, 1, 0, 34, 34, 1, 0, 39, 39, 1, 0, 10, 10, 111,
		0, 1, 1, 0, 0, 0, 0, 3, 1, 0, 0, 0, 0, 5, 1, 0, 0, 0, 0, 7, 1, 0, 0, 0,
		0, 9, 1, 0, 0, 0, 0, 11, 1, 0, 0, 0, 0, 13, 1, 0, 0, 0, 0, 15, 1, 0, 0,
		0, 0, 17, 1, 0, 0, 0, 0, 19, 1, 0, 0, 0, 1, 21, 1, 0, 0, 0, 3, 26, 1, 0,
		0, 0, 5, 49, 1, 0, 0, 0, 7, 51, 1, 0, 0, 0, 9, 53, 1, 0, 0, 0, 11, 61,
		1, 0, 0, 0, 13, 67, 1, 0, 0, 0, 15, 74, 1, 0, 0, 0, 17, 87, 1, 0, 0, 0,
		19, 98, 1, 0, 0, 0, 21, 22, 7, 0, 0, 0, 22, 23, 1, 0, 0, 0, 23, 24, 6,
		0, 0, 0, 24, 2, 1, 0, 0, 0, 25, 27, 7, 1, 0, 0, 26, 25, 1, 0, 0, 0, 27,
		28, 1, 0, 0, 0, 28, 26, 1, 0, 0, 0, 28, 29, 1, 0, 0, 0, 29, 4, 1, 0, 0,
		0, 30, 36, 5, 34, 0, 0, 31, 35, 8, 2, 0, 0, 32, 33, 5, 92, 0, 0, 33, 35,
		5, 34, 0, 0, 34, 31, 1, 0, 0, 0, 34, 32, 1, 0, 0, 0, 35, 38, 1, 0, 0, 0,
		36, 34, 1, 0, 0, 0, 36, 37, 1, 0, 0, 0, 37, 39, 1, 0, 0, 0, 38, 36, 1,
		0, 0, 0, 39, 50, 5, 34, 0, 0, 40, 46, 5, 39, 0, 0, 41, 45, 8, 3, 0, 0,
		42, 43, 5, 92, 0, 0, 43, 45, 5, 39, 0, 0, 44, 41, 1, 0, 0, 0, 44, 42, 1,
		0, 0, 0, 45, 48, 1, 0, 0, 0, 46, 44, 1, 0, 0, 0, 46, 47, 1, 0, 0, 0, 47,
		50, 1, 0, 0, 0, 48, 46, 1, 0, 0, 0, 49, 30, 1, 0, 0, 0, 49, 40, 1, 0, 0,
		0, 50, 6, 1, 0, 0, 0, 51, 52, 5, 10, 0, 0, 52, 8, 1, 0, 0, 0, 53, 54, 5,
		109, 0, 0, 54, 55, 5, 115, 0, 0, 55, 56, 5, 103, 0, 0, 56, 57, 5, 99, 0,
		0, 57, 58, 5, 116, 0, 0, 58, 59, 5, 120, 0, 0, 59, 60, 5, 116, 0, 0, 60,
		10, 1, 0, 0, 0, 61, 62, 5, 109, 0, 0, 62, 63, 5, 115, 0, 0, 63, 64, 5,
		103, 0, 0, 64, 65, 5, 105, 0, 0, 65, 66, 5, 100, 0, 0, 66, 12, 1, 0, 0,
		0, 67, 68, 5, 109, 0, 0, 68, 69, 5, 115, 0, 0, 69, 70, 5, 103, 0, 0, 70,
		71, 5, 115, 0, 0, 71, 72, 5, 116, 0, 0, 72, 73, 5, 114, 0, 0, 73, 14, 1,
		0, 0, 0, 74, 75, 5, 109, 0, 0, 75, 76, 5, 115, 0, 0, 76, 77, 5, 103, 0,
		0, 77, 78, 5, 105, 0, 0, 78, 79, 5, 100, 0, 0, 79, 80, 5, 95, 0, 0, 80,
		81, 5, 112, 0, 0, 81, 82, 5, 108, 0, 0, 82, 83, 5, 117, 0, 0, 83, 84, 5,
		114, 0, 0, 84, 85, 5, 97, 0, 0, 85, 86, 5, 108, 0, 0, 86, 16, 1, 0, 0,
		0, 87, 88, 5, 109, 0, 0, 88, 89, 5, 115, 0, 0, 89, 90, 5, 103, 0, 0, 90,
		91, 5, 115, 0, 0, 91, 92, 5, 116, 0, 0, 92, 93, 5, 114, 0, 0, 93, 94, 5,
		91, 0, 0, 94, 95, 1, 0, 0, 0, 95, 96, 3, 3, 1, 0, 96, 97, 5, 93, 0, 0,
		97, 18, 1, 0, 0, 0, 98, 102, 5, 35, 0, 0, 99, 101, 8, 4, 0, 0, 100, 99,
		1, 0, 0, 0, 101, 104, 1, 0, 0, 0, 102, 100, 1, 0, 0, 0, 102, 103, 1, 0,
		0, 0, 103, 20, 1, 0, 0, 0, 104, 102, 1, 0, 0, 0, 8, 0, 28, 34, 36, 44,
		46, 49, 102, 1, 6, 0, 0,
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

// PoLexerInit initializes any static state used to implement PoLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewPoLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func PoLexerInit() {
	staticData := &PoLexerLexerStaticData
	staticData.once.Do(polexerLexerInit)
}

// NewPoLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewPoLexer(input antlr.CharStream) *PoLexer {
	PoLexerInit()
	l := new(PoLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
	staticData := &PoLexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	l.channelNames = staticData.ChannelNames
	l.modeNames = staticData.ModeNames
	l.RuleNames = staticData.RuleNames
	l.LiteralNames = staticData.LiteralNames
	l.SymbolicNames = staticData.SymbolicNames
	l.GrammarFileName = "Po.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// PoLexer tokens.
const (
	PoLexerWS            = 1
	PoLexerINT           = 2
	PoLexerSTRING        = 3
	PoLexerNL            = 4
	PoLexerMSGCTXT       = 5
	PoLexerMSGID         = 6
	PoLexerMSGSTR        = 7
	PoLexerPLURAL_MSGID  = 8
	PoLexerPLURAL_MSGSTR = 9
	PoLexerCOMMENT       = 10
)
