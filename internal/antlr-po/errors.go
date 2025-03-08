package parse

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
)

var _ antlr.ErrorListener = (*CustomErrorListener)(nil)

type CustomErrorListener struct {
	*antlr.DefaultErrorListener          // Extends the default behavior
	Errors                      []string // Slice to store errors
}

func (c *CustomErrorListener) SyntaxError(
	_ antlr.Recognizer,
	_ any,
	line, column int,
	msg string,
	_ antlr.RecognitionException,
) {
	errorMsg := fmt.Sprintf("Syntax Error at line %d:%d - %s", line, column, msg)
	c.Errors = append(c.Errors, errorMsg)
}

func (c *CustomErrorListener) ReportAmbiguity(
	_ antlr.Parser,
	_ *antlr.DFA,
	startIndex, stopIndex int,
	_ bool,
	_ *antlr.BitSet,
	_ *antlr.ATNConfigSet,
) {
	errorMsg := fmt.Sprintf("Ambiguity detected between indices %d and %d", startIndex, stopIndex)
	c.Errors = append(c.Errors, errorMsg)
}

func (c *CustomErrorListener) ReportAttemptingFullContext(
	_ antlr.Parser,
	_ *antlr.DFA,
	startIndex, stopIndex int,
	_ *antlr.BitSet,
	_ *antlr.ATNConfigSet,
) {
	errorMsg := fmt.Sprintf(
		"Attempting full context resolution between indices %d and %d",
		startIndex,
		stopIndex,
	)
	c.Errors = append(c.Errors, errorMsg)
}

func (c *CustomErrorListener) ReportContextSensitivity(
	_ antlr.Parser,
	_ *antlr.DFA,
	startIndex, stopIndex, _ int,
	_ *antlr.ATNConfigSet,
) {
	errorMsg := fmt.Sprintf(
		"Context sensitivity detected between indices %d and %d",
		startIndex,
		stopIndex,
	)
	c.Errors = append(c.Errors, errorMsg)
}
