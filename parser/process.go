package parser

import (
	"go/ast"
	"go/token"
	"strconv"

	"github.com/gookit/color"
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
) (translation Translation, valid bool) {
	def := gotextGetter[method]

	id := extractArgument(callExpr, def.ID)
	valid = id.valid
	if id.valid {
		translation.ID = id.str
		translation.Locations = append(
			translation.Locations,
			Location{findLine(f.content, id.pos), f.path},
		)
	}
	context := extractArgument(callExpr, def.Context)
	if context.valid {
		translation.Context = context.str
	}
	plural := extractArgument(callExpr, def.Plural)
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

func extractArgument(callExpr *ast.CallExpr, index int) (a arg) {
	if index == -1 {
		return
	}
	basicLit, ok := callExpr.Args[index].(*ast.BasicLit)
	if !ok {
		return
	}
	if basicLit.Kind != token.STRING {
		return
	}

	content, err := strconv.Unquote(basicLit.Value)
	if err != nil {
		color.Errorln(err)
		return
	}

	return arg{
		str:   content,
		valid: true,
		pos:   basicLit.Pos(),
	}
}
