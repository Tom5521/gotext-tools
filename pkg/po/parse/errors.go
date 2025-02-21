package parse

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
)

var _ antlr.ErrorListener = (*CustomErrorListener)(nil)

// CustomErrorListener is a custom implementation of antlr.ErrorListener
type CustomErrorListener struct {
	*antlr.DefaultErrorListener          // Extends the default behavior
	Errors                      []string // Slice to store errors
}

// SyntaxError handles syntax errors and appends them to the Errors slice
func (c *CustomErrorListener) SyntaxError(
	recognizer antlr.Recognizer,
	offendingSymbol interface{},
	line, column int,
	msg string,
	e antlr.RecognitionException,
) {
	errorMsg := fmt.Sprintf("Syntax Error at line %d:%d - %s", line, column, msg)
	c.Errors = append(c.Errors, errorMsg)
}

// ReportAmbiguity handles grammar ambiguities and appends them to the Errors slice
func (c *CustomErrorListener) ReportAmbiguity(
	recognizer antlr.Parser,
	dfa *antlr.DFA,
	startIndex, stopIndex int,
	exact bool,
	ambigAlts *antlr.BitSet,
	configs *antlr.ATNConfigSet,
) {
	errorMsg := fmt.Sprintf("Ambiguity detected between indices %d and %d", startIndex, stopIndex)
	c.Errors = append(c.Errors, errorMsg)
}

// ReportAttemptingFullContext handles full context resolution attempts and appends them to the Errors slice
func (c *CustomErrorListener) ReportAttemptingFullContext(
	recognizer antlr.Parser,
	dfa *antlr.DFA,
	startIndex, stopIndex int,
	conflictingAlts *antlr.BitSet,
	configs *antlr.ATNConfigSet,
) {
	errorMsg := fmt.Sprintf(
		"Attempting full context resolution between indices %d and %d",
		startIndex,
		stopIndex,
	)
	c.Errors = append(c.Errors, errorMsg)
}

// ReportContextSensitivity handles context sensitivity and appends them to the Errors slice
func (c *CustomErrorListener) ReportContextSensitivity(
	recognizer antlr.Parser,
	dfa *antlr.DFA,
	startIndex, stopIndex, prediction int,
	configs *antlr.ATNConfigSet,
) {
	errorMsg := fmt.Sprintf(
		"Context sensitivity detected between indices %d and %d",
		startIndex,
		stopIndex,
	)
	c.Errors = append(c.Errors, errorMsg)
}
