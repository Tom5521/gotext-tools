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
		"PLURAL_MSGSTR", "COMMENT", "FLAG_COMMENT", "EXTRACTED_COMMENT", "REFERENCE_COMMENT",
		"PREVIOUS_COMMENT",
	}
	staticData.RuleNames = []string{
		"WS", "INT", "STRING", "NL", "MSGCTXT", "MSGID", "MSGSTR", "PLURAL_MSGID",
		"PLURAL_MSGSTR", "COMMENT", "FLAG_COMMENT", "EXTRACTED_COMMENT", "REFERENCE_COMMENT",
		"PREVIOUS_COMMENT",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 14, 149, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2,
		10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 1, 0, 1, 0, 1, 0,
		1, 0, 1, 1, 4, 1, 35, 8, 1, 11, 1, 12, 1, 36, 1, 2, 1, 2, 1, 2, 1, 2, 5,
		2, 43, 8, 2, 10, 2, 12, 2, 46, 9, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 5, 2,
		53, 8, 2, 10, 2, 12, 2, 56, 9, 2, 3, 2, 58, 8, 2, 1, 3, 1, 3, 1, 4, 1,
		4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1,
		5, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 7, 1, 7, 1, 7, 1, 7, 1,
		7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 8, 1, 8, 1, 8, 1,
		8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 9, 1, 9, 5, 9, 109, 8,
		9, 10, 9, 12, 9, 112, 9, 9, 1, 10, 1, 10, 1, 10, 1, 10, 5, 10, 118, 8,
		10, 10, 10, 12, 10, 121, 9, 10, 1, 11, 1, 11, 1, 11, 1, 11, 5, 11, 127,
		8, 11, 10, 11, 12, 11, 130, 9, 11, 1, 12, 1, 12, 1, 12, 1, 12, 5, 12, 136,
		8, 12, 10, 12, 12, 12, 139, 9, 12, 1, 13, 1, 13, 1, 13, 1, 13, 5, 13, 145,
		8, 13, 10, 13, 12, 13, 148, 9, 13, 0, 0, 14, 1, 1, 3, 2, 5, 3, 7, 4, 9,
		5, 11, 6, 13, 7, 15, 8, 17, 9, 19, 10, 21, 11, 23, 12, 25, 13, 27, 14,
		1, 0, 5, 3, 0, 9, 10, 13, 13, 32, 32, 1, 0, 48, 57, 1, 0, 34, 34, 1, 0,
		39, 39, 1, 0, 10, 10, 159, 0, 1, 1, 0, 0, 0, 0, 3, 1, 0, 0, 0, 0, 5, 1,
		0, 0, 0, 0, 7, 1, 0, 0, 0, 0, 9, 1, 0, 0, 0, 0, 11, 1, 0, 0, 0, 0, 13,
		1, 0, 0, 0, 0, 15, 1, 0, 0, 0, 0, 17, 1, 0, 0, 0, 0, 19, 1, 0, 0, 0, 0,
		21, 1, 0, 0, 0, 0, 23, 1, 0, 0, 0, 0, 25, 1, 0, 0, 0, 0, 27, 1, 0, 0, 0,
		1, 29, 1, 0, 0, 0, 3, 34, 1, 0, 0, 0, 5, 57, 1, 0, 0, 0, 7, 59, 1, 0, 0,
		0, 9, 61, 1, 0, 0, 0, 11, 69, 1, 0, 0, 0, 13, 75, 1, 0, 0, 0, 15, 82, 1,
		0, 0, 0, 17, 95, 1, 0, 0, 0, 19, 106, 1, 0, 0, 0, 21, 113, 1, 0, 0, 0,
		23, 122, 1, 0, 0, 0, 25, 131, 1, 0, 0, 0, 27, 140, 1, 0, 0, 0, 29, 30,
		7, 0, 0, 0, 30, 31, 1, 0, 0, 0, 31, 32, 6, 0, 0, 0, 32, 2, 1, 0, 0, 0,
		33, 35, 7, 1, 0, 0, 34, 33, 1, 0, 0, 0, 35, 36, 1, 0, 0, 0, 36, 34, 1,
		0, 0, 0, 36, 37, 1, 0, 0, 0, 37, 4, 1, 0, 0, 0, 38, 44, 5, 34, 0, 0, 39,
		43, 8, 2, 0, 0, 40, 41, 5, 92, 0, 0, 41, 43, 5, 34, 0, 0, 42, 39, 1, 0,
		0, 0, 42, 40, 1, 0, 0, 0, 43, 46, 1, 0, 0, 0, 44, 42, 1, 0, 0, 0, 44, 45,
		1, 0, 0, 0, 45, 47, 1, 0, 0, 0, 46, 44, 1, 0, 0, 0, 47, 58, 5, 34, 0, 0,
		48, 54, 5, 39, 0, 0, 49, 53, 8, 3, 0, 0, 50, 51, 5, 92, 0, 0, 51, 53, 5,
		39, 0, 0, 52, 49, 1, 0, 0, 0, 52, 50, 1, 0, 0, 0, 53, 56, 1, 0, 0, 0, 54,
		52, 1, 0, 0, 0, 54, 55, 1, 0, 0, 0, 55, 58, 1, 0, 0, 0, 56, 54, 1, 0, 0,
		0, 57, 38, 1, 0, 0, 0, 57, 48, 1, 0, 0, 0, 58, 6, 1, 0, 0, 0, 59, 60, 5,
		10, 0, 0, 60, 8, 1, 0, 0, 0, 61, 62, 5, 109, 0, 0, 62, 63, 5, 115, 0, 0,
		63, 64, 5, 103, 0, 0, 64, 65, 5, 99, 0, 0, 65, 66, 5, 116, 0, 0, 66, 67,
		5, 120, 0, 0, 67, 68, 5, 116, 0, 0, 68, 10, 1, 0, 0, 0, 69, 70, 5, 109,
		0, 0, 70, 71, 5, 115, 0, 0, 71, 72, 5, 103, 0, 0, 72, 73, 5, 105, 0, 0,
		73, 74, 5, 100, 0, 0, 74, 12, 1, 0, 0, 0, 75, 76, 5, 109, 0, 0, 76, 77,
		5, 115, 0, 0, 77, 78, 5, 103, 0, 0, 78, 79, 5, 115, 0, 0, 79, 80, 5, 116,
		0, 0, 80, 81, 5, 114, 0, 0, 81, 14, 1, 0, 0, 0, 82, 83, 5, 109, 0, 0, 83,
		84, 5, 115, 0, 0, 84, 85, 5, 103, 0, 0, 85, 86, 5, 105, 0, 0, 86, 87, 5,
		100, 0, 0, 87, 88, 5, 95, 0, 0, 88, 89, 5, 112, 0, 0, 89, 90, 5, 108, 0,
		0, 90, 91, 5, 117, 0, 0, 91, 92, 5, 114, 0, 0, 92, 93, 5, 97, 0, 0, 93,
		94, 5, 108, 0, 0, 94, 16, 1, 0, 0, 0, 95, 96, 5, 109, 0, 0, 96, 97, 5,
		115, 0, 0, 97, 98, 5, 103, 0, 0, 98, 99, 5, 115, 0, 0, 99, 100, 5, 116,
		0, 0, 100, 101, 5, 114, 0, 0, 101, 102, 5, 91, 0, 0, 102, 103, 1, 0, 0,
		0, 103, 104, 3, 3, 1, 0, 104, 105, 5, 93, 0, 0, 105, 18, 1, 0, 0, 0, 106,
		110, 5, 35, 0, 0, 107, 109, 8, 4, 0, 0, 108, 107, 1, 0, 0, 0, 109, 112,
		1, 0, 0, 0, 110, 108, 1, 0, 0, 0, 110, 111, 1, 0, 0, 0, 111, 20, 1, 0,
		0, 0, 112, 110, 1, 0, 0, 0, 113, 114, 5, 35, 0, 0, 114, 115, 5, 44, 0,
		0, 115, 119, 1, 0, 0, 0, 116, 118, 8, 4, 0, 0, 117, 116, 1, 0, 0, 0, 118,
		121, 1, 0, 0, 0, 119, 117, 1, 0, 0, 0, 119, 120, 1, 0, 0, 0, 120, 22, 1,
		0, 0, 0, 121, 119, 1, 0, 0, 0, 122, 123, 5, 35, 0, 0, 123, 124, 5, 46,
		0, 0, 124, 128, 1, 0, 0, 0, 125, 127, 8, 4, 0, 0, 126, 125, 1, 0, 0, 0,
		127, 130, 1, 0, 0, 0, 128, 126, 1, 0, 0, 0, 128, 129, 1, 0, 0, 0, 129,
		24, 1, 0, 0, 0, 130, 128, 1, 0, 0, 0, 131, 132, 5, 35, 0, 0, 132, 133,
		5, 58, 0, 0, 133, 137, 1, 0, 0, 0, 134, 136, 8, 4, 0, 0, 135, 134, 1, 0,
		0, 0, 136, 139, 1, 0, 0, 0, 137, 135, 1, 0, 0, 0, 137, 138, 1, 0, 0, 0,
		138, 26, 1, 0, 0, 0, 139, 137, 1, 0, 0, 0, 140, 141, 5, 35, 0, 0, 141,
		142, 5, 124, 0, 0, 142, 146, 1, 0, 0, 0, 143, 145, 8, 4, 0, 0, 144, 143,
		1, 0, 0, 0, 145, 148, 1, 0, 0, 0, 146, 144, 1, 0, 0, 0, 146, 147, 1, 0,
		0, 0, 147, 28, 1, 0, 0, 0, 148, 146, 1, 0, 0, 0, 12, 0, 36, 42, 44, 52,
		54, 57, 110, 119, 128, 137, 146, 1, 6, 0, 0,
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
	PoLexerWS                = 1
	PoLexerINT               = 2
	PoLexerSTRING            = 3
	PoLexerNL                = 4
	PoLexerMSGCTXT           = 5
	PoLexerMSGID             = 6
	PoLexerMSGSTR            = 7
	PoLexerPLURAL_MSGID      = 8
	PoLexerPLURAL_MSGSTR     = 9
	PoLexerCOMMENT           = 10
	PoLexerFLAG_COMMENT      = 11
	PoLexerEXTRACTED_COMMENT = 12
	PoLexerREFERENCE_COMMENT = 13
	PoLexerPREVIOUS_COMMENT  = 14
)
