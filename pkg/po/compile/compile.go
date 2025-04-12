package compile

import (
	"io"

	"github.com/Tom5521/gotext-tools/pkg/po"
)

type eOf interface{ *po.File | po.Entries }

func file[T eOf](i T) *po.File {
	var f *po.File
	switch v := any(i).(type) {
	case *po.File:
		f = v
	case po.Entries:
		f = &po.File{Entries: v}
	}

	return f
}

func PoToWriter[T eOf](f T, w io.Writer, opts ...PoOption) error {
	return NewPo(file(f), opts...).ToWriter(w)
}

func PoToString[T eOf](f T, opts ...PoOption) string {
	return NewPo(file(f), opts...).ToString()
}

func PoToFile[T eOf](f T, path string, opts ...PoOption) error {
	return NewPo(file(f), opts...).ToFile(path)
}

func PoToBytes[T eOf](f T, opts ...PoOption) []byte {
	return NewPo(file(f), opts...).ToBytes()
}

func MoToWriter[T eOf](f T, w io.Writer, opts ...MoOption) error {
	return NewMo(file(f), opts...).ToWriter(w)
}

func MoToBytes[T eOf](f T, opts ...MoOption) []byte {
	return NewMo(file(f), opts...).ToBytes()
}

func MoToFile[T eOf](f T, path string, opts ...MoOption) error {
	return NewMo(file(f), opts...).ToFile(path)
}
