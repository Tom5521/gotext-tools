package parser

import (
	"go/ast"
	"go/token"
)

func (f *File) processMethod(method string, callExpr *ast.CallExpr) (Translation, bool) {
	switch method {
	case Get, GetD:
		return handleSimpleGet(f.path, f.content, methodToIndex(method), callExpr)
	case GetN:
		return f.processGetN(callExpr)
	case GetC:
		return f.processGetC(callExpr)
	case GetND:
		return f.processGetND(callExpr)
	case GetNC:
		return f.processGetNC(callExpr)
	case GetNDC:
		return f.processGetNDC(callExpr)
	default:
		return Translation{}, false
	}
}

func handleSimpleGet(path, content string, index int, callExpr *ast.CallExpr) (Translation, bool) {
	id, pos, valid := extractArgument(callExpr, index)
	if !valid {
		return Translation{}, false
	}
	return Translation{
		ID: id,
		Locations: []Location{
			{File: path, Line: findLine(content, pos)},
		},
	}, true
}

func extractArgument(callExpr *ast.CallExpr, index int) (string, token.Pos, bool) {
	basicLit, ok := callExpr.Args[index].(*ast.BasicLit)
	if !ok {
		return "", 0, false
	}
	content := basicLit.Value[1 : len(basicLit.Value)-1] // Strip surrounding quotes.
	return content, basicLit.ValuePos, true
}

func (f *File) processGetN(callExpr *ast.CallExpr) (Translation, bool) {
	translation, valid := handleSimpleGet(f.path, f.content, 0, callExpr)
	if !valid {
		return Translation{}, false
	}

	plural, _, valid := extractArgument(callExpr, 1)
	if !valid {
		return Translation{}, false
	}
	translation.Plural = plural
	return translation, true
}

func (f *File) processGetC(callExpr *ast.CallExpr) (Translation, bool) {
	translation, valid := handleSimpleGet(f.path, f.content, 0, callExpr)
	if !valid {
		return Translation{}, false
	}

	context, _, valid := extractArgument(callExpr, 1)
	if !valid {
		return Translation{}, false
	}

	translation.Context = context

	return translation, true
}

func (f *File) processGetND(callExpr *ast.CallExpr) (Translation, bool) {
	translation, valid := handleSimpleGet(f.path, f.content, 1, callExpr)
	if !valid {
		return Translation{}, false
	}

	plural, _, valid := extractArgument(callExpr, 2)
	if !valid {
		return Translation{}, false
	}
	translation.Plural = plural

	return translation, true
}

func (f *File) processGetNC(callExpr *ast.CallExpr) (Translation, bool) {
	translation, valid := handleSimpleGet(f.path, f.content, 0, callExpr)
	if !valid {
		return Translation{}, false
	}

	plural, _, valid := extractArgument(callExpr, 1)
	if !valid {
		return Translation{}, false
	}
	translation.Plural = plural

	context, _, valid := extractArgument(callExpr, 3)
	if !valid {
		return Translation{}, false
	}

	translation.Context = context

	return translation, true
}

func (f *File) processGetNDC(callExpr *ast.CallExpr) (Translation, bool) {
	translation, valid := handleSimpleGet(f.path, f.content, 1, callExpr)
	if !valid {
		return Translation{}, false
	}

	plural, _, valid := extractArgument(callExpr, 2)
	if !valid {
		return Translation{}, false
	}
	translation.Plural = plural

	context, _, valid := extractArgument(callExpr, 4)
	if !valid {
		return Translation{}, false
	}

	translation.Context = context

	return translation, true
}
