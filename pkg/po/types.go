package po

import (
	"io"
)

type Parser interface {
	Parse() *File
	Error() error
	Errors() []error
}

type Compiler interface {
	SetFile(*File)
	ToWriter(io.Writer) error
	ToBytes() []byte
}
