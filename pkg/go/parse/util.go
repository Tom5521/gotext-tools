package parse

import (
	"go/ast"
)

// InspectNode wraps the `ast.Inspect` function, allowing for convenient traversal
// of an abstract syntax tree (AST). This function returns a higher-order function
// that accepts a visitor function, which is applied to each node in the AST.
//
// Parameters:
//   - root: The root node of the AST to be traversed.
//
// Returns:
//   - A function that accepts a visitor function (func(ast.Node) bool).
//     The visitor function is called for each node in the AST.
func InspectNode(root ast.Node) func(func(ast.Node) bool) {
	return func(f func(ast.Node) bool) {
		ast.Inspect(root, f)
	}
}
