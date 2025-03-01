package compiler

import "io"

type Compiler interface {
	ToWriter(io.Writer, ...PoOption) error
	ToFile(string, ...PoOption) error
	ToString(...PoOption) string
	ToBytes(...PoOption) []byte
}
