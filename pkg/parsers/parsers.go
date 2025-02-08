package parsers

import (
	"github.com/Tom5521/xgotext/pkg/po/types"
)

type Parser interface {
	Parse() *types.File
	Errors() []error
	Warnings() []string
}

type Config struct {
	Language   string
	Exclude    []string
	ExtractAll bool
}
