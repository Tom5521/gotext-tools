package generator

import (
	"github.com/Tom5521/xgotext/pkg/po/parse/ast"
	"github.com/Tom5521/xgotext/pkg/po/types"
)

type Generator struct {
	file *ast.File
}

func Gen(input *ast.File) (output *types.File, warnings []string, errs []error)
func (g *Generator) genHeader(str string) (header types.Header)

func (g *Generator) genEntries(
	nodes []ast.Node,
) (entries []types.Entry, warns []string, errs []error)
