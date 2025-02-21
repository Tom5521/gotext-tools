package parse

import (
	"go/ast"
)

// InspectNode wraps the `ast.Inspect` function to provide convenient traversal of an abstract syntax tree (AST).
//
// ### Parameters:
// - `root`: The root node of the AST to be traversed.
//
// ### Returns:
// - A higher-order function that accepts a visitor function (`func(ast.Node) bool`).
// - The visitor function is called for each node in the AST.
//
// This function simplifies AST traversal by allowing users to focus on node-specific processing logic.
func InspectNode(root ast.Node) chan ast.Node {
	ch := make(chan ast.Node)

	go func() {
		defer close(ch)
		ast.Inspect(root, func(n ast.Node) bool {
			ch <- n
			return true
		})
	}()

	return ch
}
