package parse

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
)

var _ antlr.ErrorListener = (*CustomErrorListener)(nil)

type CustomErrorListener struct {
	*antlr.DefaultErrorListener
	Errors []string
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
