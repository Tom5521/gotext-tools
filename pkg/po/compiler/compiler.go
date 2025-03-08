package compiler

import "io"

type Compiler interface {
	ToWriter(io.Writer) error
	ToFile(string) error
	ToBytes() []byte
}
