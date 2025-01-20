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
	ID      int // Position of message ID argument
	Plural  int // Position of plural form argument (-1 if not applicable)
	Context int // Position of context argument (-1 if not applicable)
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

// isGotextCall checks if an AST node represents a gotext function call.
func (f *File) isGotextCall(n ast.Node) bool {
	callExpr, ok := n.(*ast.CallExpr)
	if !ok {
		return false
	}

	selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident, ok := selectorExpr.X.(*ast.Ident)
	if !ok || ident.Name != f.pkgName {
		return false
	}

	_, ok = translationMethods[selectorExpr.Sel.Name]
	return ok
}

func (f *File) basicLitToTranslation(n *ast.BasicLit) (entry.Translation, error) {
	str, err := strconv.Unquote(n.Value)
	if err != nil {
		return entry.Translation{}, err
	}

	return entry.Translation{
		ID: str,
		Locations: []entry.Location{{
			Line: util.FindLine(f.content, n.Pos()),
			File: f.path,
		}},
	}, nil
}

func (f *File) processGeneric(exprs ...ast.Expr) ([]entry.Translation, []error) {
	var translations []entry.Translation
	var errors []error

	for _, expr := range exprs {
		if lit, ok := expr.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			if f.seenTokens[lit] {
				continue
			}

			if lit.Value == `""` {
				continue
			}

			translation, err := f.basicLitToTranslation(lit)
			if err != nil {
				errors = append(errors, err)
				continue
			}

			translations = append(translations, translation)
			f.seenTokens[lit] = true
		}
	}

	return translations, errors
}

type argumentData struct {
	str   string
	valid bool
	err   error
	pos   token.Pos
}

func (f *File) extractArg(index int, call *ast.CallExpr) (a argumentData) {
	if index == -1 {
		return
	}
	if index < 0 || index >= len(call.Args) {
		a.err = fmt.Errorf("index (%d) out of range", index)
		return
	}
	lit, ok := call.Args[index].(*ast.BasicLit)
	if !ok {
		return
	}

	if lit.Kind != token.STRING {
		a.err = fmt.Errorf("the specified argument (%d) is not a string", index)
		return
	}

	if lit.Value == `""` {
		return
	}

	pos := lit.Pos()

	f.seenTokens[lit] = true

	str, err := strconv.Unquote(lit.Value)
	if err != nil {
		a.err = fmt.Errorf("error unquoting string: %w", err)
		return
	}

	return argumentData{str, true, err, pos}
}

func (f *File) processPoCall(
	call *ast.CallExpr,
) (translation entry.Translation, valid bool, err error) {
	selector := call.Fun.(*ast.SelectorExpr)
	method := translationMethods[selector.Sel.Name]

	args := []argumentData{
		f.extractArg(method.ID, call),
		f.extractArg(method.Context, call),
		f.extractArg(method.Plural, call),
	}

	for i, arg := range args {
		if arg.err != nil {
			err = arg.err
			return
		}
		switch i {
		case 0:
			valid = arg.valid
			translation.ID = arg.str
			translation.Locations = append(translation.Locations,
				entry.Location{
					File: f.path,
					Line: util.FindLine(f.content, arg.pos),
				},
			)
			fallthrough
		case 1:
			translation.Context = arg.str
			fallthrough
		case 2:
			translation.Plural = arg.str
		}
	}

	return
}

func (f *File) processNode(n ast.Node) ([]entry.Translation, []error) {
	if n == nil {
		return nil, nil
	}
	var translations []entry.Translation
	var errors []error

	processGeneric := func(exprs ...ast.Expr) {
		t, e := f.processGeneric(exprs...)
		translations = append(translations, t...)
		errors = append(errors, e...)
	}

	processPoCall := func(call *ast.CallExpr) {
		t, valid, err := f.processPoCall(call)
		if !valid {
			return
		}
		translations = append(translations, t)
		if err != nil {
			errors = append(errors, err)
		}
	}

	if !f.config.ExtractAll {
		if f.isGotextCall(n) {
			call := n.(*ast.CallExpr)
			processPoCall(call)
		}

		return translations, errors
	}

	switch t := n.(type) {
	case *ast.CallExpr:
		if f.isGotextCall(t) {
			processPoCall(t)
		} else {
			processGeneric(t.Args...)
		}
	case *ast.AssignStmt:
		processGeneric(t.Rhs...)
	case *ast.ValueSpec:
		processGeneric(t.Values...)
	case *ast.ReturnStmt:
		processGeneric(t.Results...)
	case *ast.KeyValueExpr:
		processGeneric(t.Value)
	case *ast.SendStmt:
		processGeneric(t.Value)
	case *ast.CompositeLit:
		processGeneric(t.Elts...)
	case *ast.BinaryExpr:
		processGeneric(t.X, t.Y)
	// Switch expressions.
	case *ast.SwitchStmt:
		processGeneric(t.Tag)
	case *ast.CaseClause:
		processGeneric(t.List...)
	}

	return translations, errors
}
