package po

import (
	"io"
)

type (
	Cmp[X any] func(a, b X) int
	Parser     interface {
		Parse() *File
		Error() error
		Errors() []error
	}

	Compiler interface {
		SetFile(*File)
		ToWriter(io.Writer) error
		ToBytes() []byte
	}
)
