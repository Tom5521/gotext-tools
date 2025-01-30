package parse

import (
	"fmt"

	"github.com/Tom5521/xgotext/pkg/po/config"
)

func validateConfig(cfg config.Config) error {
	cfgErrs := cfg.Validate()
	if len(cfgErrs) > 0 {
		return fmt.Errorf("configuration is invalid: %w", cfgErrs[0])
	}
	return nil
}
