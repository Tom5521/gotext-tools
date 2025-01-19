package goparse

import (
	"go/ast"
)

func InspectNode(root ast.Node) func(func(ast.Node) bool) {
	return func(f func(ast.Node) bool) {
		ast.Inspect(root, f)
	}
}
