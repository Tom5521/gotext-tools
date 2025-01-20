package goparse

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po/entry"
)

// translationMethod defines the structure for different getter methods.
type translationMethod struct {
	IDIndex      int // Position of message ID argument
	PluralIndex  int // Position of plural form argument (-1 if not applicable)
	ContextIndex int // Position of context argument (-1 if not applicable)
}

// Define supported translation methods.
var translationMethods = map[string]translationMethod{
	"Get":    {0, -1, -1},
	"GetN":   {0, 1, -1},
	"GetD":   {1, -1, -1},
	"GetND":  {1, 2, -1},
	"GetC":   {0, -1, 1},
	"GetNC":  {0, 1, 3},
	"GetDC":  {1, -1, 2},
	"GetNDC": {1, 2, 4},
}

// translationArgument represents a parsed translation argument.
type translationArgument struct {
	Value    string
	IsValid  bool
	Position token.Pos
}

// ProcessMethod processes a translation method call and returns the corresponding Translation.
func (f *File) processMethod(
	methodName string,
	callExpr *ast.CallExpr,
) (entry.Translation, bool, error) {
	method, exists := translationMethods[methodName]
	if !exists {
		return entry.Translation{}, false, fmt.Errorf(
			"unsupported translation method: %s",
			methodName,
		)
	}

	translation := entry.Translation{}

	// Extract message ID
	msgID, err := f.extractArgument(callExpr, method.IDIndex)
	if err != nil {
		return translation, false, fmt.Errorf("failed to extract message ID: %w", err)
	}

	if !msgID.IsValid {
		return translation, false, nil
	}

	// Set message ID and location
	translation.ID = msgID.Value
	translation.Locations = []entry.Location{{
		Line: util.FindLine(f.content, msgID.Position),
		File: f.path,
	}}

	// Extract and set context if applicable
	if method.ContextIndex >= 0 {
		if context, err := f.extractArgument(callExpr, method.ContextIndex); err == nil &&
			context.IsValid {
			translation.Context = context.Value
		} else if err != nil {
			return translation, false, fmt.Errorf("failed to extract context: %w", err)
		}
	}

	// Extract and set plural form if applicable
	if method.PluralIndex >= 0 {
		if plural, err := f.extractArgument(callExpr, method.PluralIndex); err == nil &&
			plural.IsValid {
			translation.Plural = plural.Value
		} else if err != nil {
			return translation, false, fmt.Errorf("failed to extract plural form: %w", err)
		}
	}

	return translation, true, nil
}

// extractArgument extracts and validates a string argument from the call expression.
func (f *File) extractArgument(callExpr *ast.CallExpr, index int) (translationArgument, error) {
	if index < 0 || index >= len(callExpr.Args) {
		return translationArgument{IsValid: false}, nil
	}

	arg, ok := callExpr.Args[index].(*ast.BasicLit)
	if !ok || arg.Kind != token.STRING {
		return translationArgument{IsValid: false}, nil
	}

	f.seenTokens[arg] = true

	value, err := strconv.Unquote(arg.Value)
	if err != nil {
		return translationArgument{}, fmt.Errorf("failed to unquote argument value: %w", err)
	}

	return translationArgument{
		Value:    value,
		IsValid:  true,
		Position: arg.Pos(),
	}, nil
}
