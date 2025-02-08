package parsers

import (
	"errors"

	"github.com/Tom5521/xgotext/pkg/po/types"
)

type Parser interface {
	Parse() *types.File
	Errors() []error
	Warnings() []string
}

type Config struct {
	Language   string
	Nplurals   uint
	Exclude    []string
	ExtractAll bool
}

func (s Config) Validate() error {
	if s.Nplurals == 0 {
		return errors.New("nplurals is equal to 0")
	}

	return nil
}
