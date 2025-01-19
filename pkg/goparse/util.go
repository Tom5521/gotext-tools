package goparse

import (
	"fmt"
	"go/ast"

	"github.com/Tom5521/xgotext/pkg/poconfig"
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

// validateConfig checks if a `poconfig.Config` instance is valid by calling its
// `Validate` method. If the configuration contains errors, the first error is returned.
//
// Parameters:
//   - cfg: The configuration object to validate.
//
// Returns:
//   - nil if the configuration is valid.
//   - An error describing the first validation issue, if any.
//
// Example:
//
//	err := validateConfig(cfg)
//	if err != nil {
//	    fmt.Println("Invalid configuration:", err)
//	}
func validateConfig(cfg poconfig.Config) error {
	cfgErrs := cfg.Validate()
	if len(cfgErrs) > 0 {
		return fmt.Errorf("configuration is invalid: %w", cfgErrs[0])
	}
	return nil
}

