package compiler

import "io"

type Compiler interface {
	ToWriter(io.Writer, ...Option) error
	ToFile(string, ...Option) error
	ToString(...Option) string
	ToBytes(...Option) []byte
}
