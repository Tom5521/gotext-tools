package parse

import (
	"fmt"

	"github.com/Tom5521/xgotext/pkg/parsers"
)

func validateConfig(cfg parsers.Config) error {
	err := cfg.Validate()
	if err != nil {
		return fmt.Errorf("configuration is invalid: %w", err)
	}
	return nil
}
