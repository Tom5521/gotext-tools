package compiler

import "io"

type Compiler[Option MoOption | PoOption] interface {
	ToWriter(w io.Writer, options ...Option) error
	ToBytes(options ...Option) []byte
	ToFile(file string, options ...Option) error
}
