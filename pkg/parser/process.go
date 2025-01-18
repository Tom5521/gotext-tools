package parser

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
)

type getterDef struct {
	ID      int
	Plural  int
	Context int
}

var gotextGetter = map[string]getterDef{
	"Get":    {0, -1, -1},
	"GetN":   {0, 1, -1},
	"GetD":   {1, -1, -1},
	"GetND":  {1, 2, -1},
	"GetC":   {0, -1, 1},
	"GetNC":  {0, 1, 3},
	"GetDC":  {1, -1, 2},
	"GetNDC": {1, 2, 4},
}

func (f *File) processMethod(
	method string,
	callExpr *ast.CallExpr,
) (translation Translation, valid bool, err error) {
	def := gotextGetter[method]

	id, err := extractArgument(callExpr, def.ID)
	valid = id.valid && err == nil
	if err != nil {
		err = fmt.Errorf("error extracting msgid: %w", err)
		return
	}
	if id.valid {
		translation.ID = id.str
		translation.Locations = append(
			translation.Locations,
			Location{findLine(f.content, id.pos), f.path},
		)
	}
	context, err := extractArgument(callExpr, def.Context)
	if err != nil {
		err = fmt.Errorf("error extracting context: %w", err)
		return
	}
	if context.valid {
		translation.Context = context.str
	}
	plural, err := extractArgument(callExpr, def.Plural)
	if err != nil {
		err = fmt.Errorf("error extracting plural: %w", err)
		return
	}
	if plural.valid {
		translation.Plural = plural.str
	}

	return
}

type arg struct {
	str   string
	valid bool
	pos   token.Pos
}

func extractArgument(callExpr *ast.CallExpr, index int) (arg, error) {
	if index == -1 {
		return arg{}, nil
	}
	basicLit, ok := callExpr.Args[index].(*ast.BasicLit)
	if !ok {
		return arg{}, nil
	}
	if basicLit.Kind != token.STRING {
		return arg{}, fmt.Errorf("unexpected type: %v, expected: %v", token.STRING, basicLit.Kind)
	}

	content, err := strconv.Unquote(basicLit.Value)
	if err != nil {
		return arg{}, fmt.Errorf("error unquoting the argument value: %w", err)
	}

	return arg{
		str:   content,
		valid: true,
		pos:   basicLit.Pos(),
	}, nil
}
