package parse

import (
	"go/ast"
	"go/token"
	"strconv"

	"github.com/Tom5521/gotext-tools/v2/internal/util"
	"github.com/Tom5521/gotext-tools/v2/pkg/po"
)

// translationMethod defines the structure for different getter methods.
// It contains the positions of the message ID, plural form, and context arguments.
type translationMethod struct {
	ID      int // Position of message ID argument
	Plural  int // Position of plural form argument (-1 if not applicable)
	Context int // Position of context argument (-1 if not applicable)
}

// translationMethods maps method names to their respective argument positions.
var translationMethods = map[string]translationMethod{
	"Get":   {0, -1, -1}, // (str string, vars ...interface{})
	"GetN":  {0, 1, -1},  // (str string, plural string, n int, vars ...interface{})
	"GetD":  {1, -1, -1}, // (dom string, str string, vars ...interface{})
	"GetND": {1, 2, -1},  // (dom string, str string, plural string, n int, vars ...interface{})
	"GetC":  {0, -1, 1},  // (str string, ctx string, vars ...interface{})
	"GetNC": {0, 1, 3},   // (str string, plural string, n int, ctx string, vars ...interface{})
	"GetDC": {1, -1, 2},  // (dom string, str string, ctx string, vars ...interface{})
	"GetNDC": {
		1,
		2,
		4,
	}, // (dom string, str string, plural string, n int, ctx string, vars ...interface{})
}

// isGotextCall checks if an AST node represents a gotext GET function call.
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

// basicLitToEntry converts a basic literal AST node to a translation entry.
func (f *File) basicLitToEntry(n *ast.BasicLit) (po.Entry, error) {
	str, err := strconv.Unquote(n.Value)
	if err != nil {
		return po.Entry{}, f.error("error unquoting basic literal: %w", err)
	}

	return po.Entry{
		ID: str,
		Locations: []po.Location{{
			Line: util.FindLineFromReader(f.reader, n.Pos()),
			File: f.name,
		}},
	}, nil
}

// argumentData holds information about an argument extracted from a function call.
type argumentData struct {
	str   string
	valid bool
	err   error
	pos   token.Pos
}

// extractArg extracts a string argument from a function call at the specified index.
func (f *File) extractArg(index int, call *ast.CallExpr) (a argumentData) {
	if index == -1 {
		return
	}
	if index < 0 || index >= len(call.Args) {
		a.err = f.error("index (%d) out of range", index)
		return
	}
	lit, ok := call.Args[index].(*ast.BasicLit)
	if !ok {
		return
	}

	if lit.Kind != token.STRING {
		a.err = f.error("the specified argument (%d) is not a string", index)
		return
	}

	if lit.Value == `""` {
		return
	}

	f.seenNodes[lit] = struct{}{}

	str, err := strconv.Unquote(lit.Value)
	if err != nil {
		a.err = f.error("error unquoting basic literal: %w", err)
		return
	}

	return argumentData{str, true, err, lit.Pos()}
}

// processPoCall processes a gotext function call and extracts translation entries.
func (f *File) processPoCall(
	call *ast.CallExpr,
) (entry po.Entry, valid bool, err error) {
	selector, _ := call.Fun.(*ast.SelectorExpr)
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
			entry.ID = arg.str
			entry.Locations = append(entry.Locations,
				po.Location{
					File: f.name,
					Line: util.FindLineFromReader(f.reader, arg.pos),
				},
			)
			fallthrough
		case 1:
			entry.Context = arg.str
			fallthrough
		case 2:
			entry.Plural = arg.str
		}
	}

	return
}

// processNode processes an AST node and extracts translation entries.
func (f *File) processNode(n ast.Node) (po.Entries, []error) {
	if n == nil {
		return nil, nil
	}
	var entries po.Entries
	var errors []error

	processPoCall := func(call *ast.CallExpr) {
		t, valid, err := f.processPoCall(call)
		if err != nil || !valid {
			return
		}

		entries = append(entries, t)
	}

	if !f.config.ExtractAll {
		if f.isGotextCall(n) {
			call, _ := n.(*ast.CallExpr)
			processPoCall(call)
		}

		return entries, errors
	}

	switch t := n.(type) {
	case *ast.ImportSpec:
		f.seenNodes[t.Path] = struct{}{}
	case *ast.CallExpr:
		if f.isGotextCall(t) {
			processPoCall(t)
		}
	case *ast.BasicLit:
		_, ok := f.seenNodes[t]
		if t.Kind != token.STRING || ok || t.Value == `""` {
			break
		}
		f.seenNodes[t] = struct{}{}
		entry, err := f.basicLitToEntry(t)
		if err != nil {
			errors = append(errors, err)
			break
		}
		entries = append(entries, entry)
	}

	return entries, errors
}
