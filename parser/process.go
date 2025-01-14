package parser

import (
	"go/ast"
	"go/token"
	"strconv"

	"github.com/gookit/color"
)

type GetterDef struct {
	ID      int
	Plural  int
	Context int
}

var gotextGetter = map[string]GetterDef{
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

	id, pos, ok := extractArgument(callExpr, def.ID)
	if ok {
		valid = true
		translation.ID = id
		translation.Locations = append(
			translation.Locations,
			Location{findLine(f.content, pos), f.path},
		)
	}

	context, _, ok := extractArgument(callExpr, def.Context)
	if ok {
		translation.Context = context
	}

	plural, _, ok := extractArgument(callExpr, def.Plural)
	if ok {
		translation.Plural = plural
	}

	return
}

func extractArgument(callExpr *ast.CallExpr, index int) (string, token.Pos, bool) {
	if index == -1 {
		return "", 0, false
	}
	basicLit, ok := callExpr.Args[index].(*ast.BasicLit)
	if !ok {
		return "", 0, false
	}
	content, err := strconv.Unquote(basicLit.Value)
	if err != nil {
		color.Errorln(err)
		return "", 0, false
	}
	return content, basicLit.ValuePos, true
}
