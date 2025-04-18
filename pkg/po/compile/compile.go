package compile

import (
	"io"

	"github.com/Tom5521/gotext-tools/v2/pkg/po"
)

func file[T po.EntriesOrFile](i T) *po.File {
	var f *po.File
	switch v := any(i).(type) {
	case *po.File:
		f = v
	case po.Entries:
		f = &po.File{Entries: v}
	}

	return f
}

func PoToWriter[T po.EntriesOrFile](f T, w io.Writer, opts ...PoOption) error {
	return NewPo(file(f), opts...).ToWriter(w)
}

func PoToString[T po.EntriesOrFile](f T, opts ...PoOption) string {
	return NewPo(file(f), opts...).ToString()
}

func PoToFile[T po.EntriesOrFile](f T, path string, opts ...PoOption) error {
	return NewPo(file(f), opts...).ToFile(path)
}

func PoToBytes[T po.EntriesOrFile](f T, opts ...PoOption) []byte {
	return NewPo(file(f), opts...).ToBytes()
}

func MoToWriter[T po.EntriesOrFile](f T, w io.Writer, opts ...MoOption) error {
	return NewMo(file(f), opts...).ToWriter(w)
}

func MoToBytes[T po.EntriesOrFile](f T, opts ...MoOption) []byte {
	return NewMo(file(f), opts...).ToBytes()
}

func MoToFile[T po.EntriesOrFile](f T, path string, opts ...MoOption) error {
	return NewMo(file(f), opts...).ToFile(path)
}
