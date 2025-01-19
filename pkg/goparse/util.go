package goparse

import (
	"fmt"
	"go/ast"

	"github.com/Tom5521/xgotext/pkg/poconfig"
)

func InspectNode(root ast.Node) func(func(ast.Node) bool) {
	return func(f func(ast.Node) bool) {
		ast.Inspect(root, f)
	}
}

func validateConfig(cfg poconfig.Config) error {
	cfgErrs := cfg.Validate()
	if len(cfgErrs) > 0 {
		return fmt.Errorf("configuration is invalid: %w", cfgErrs[0])
	}
	return nil
}
