package parse

import (
	"io"
	"os"

	"github.com/Tom5521/gotext-tools/v2/pkg/po"
)

// Parse directly the provided file.
func Po(path string, opts ...PoOption) (*po.File, error) {
	parser, err := NewPo(path, opts...)
	if err != nil {
		return nil, err
	}
	file := parser.Parse()

	return file, parser.Error()
}

func PoFromReader(r io.Reader, name string, opts ...PoOption) (*po.File, error) {
	parser, err := NewPoFromReader(r, name, opts...)
	if err != nil {
		return nil, err
	}

	file := parser.Parse()
	return file, parser.Error()
}

func PoFromFile(f *os.File, opts ...PoOption) (*po.File, error) {
	parser, err := NewPoFromFile(f, opts...)
	if err != nil {
		return nil, err
	}
	file := parser.Parse()
	return file, parser.Error()
}

func PoFromString(s, name string, opts ...PoOption) (*po.File, error) {
	parser := NewPoFromString(s, name, opts...)

	file := parser.Parse()
	return file, parser.Error()
}

func PoFromBytes(b []byte, name string, opts ...PoOption) (*po.File, error) {
	parser := NewPoFromBytes(b, name, opts...)
	file := parser.Parse()

	return file, parser.Error()
}

func Mo(path string) (*po.File, error) {
	parser, err := NewMo(path)
	if err != nil {
		return nil, err
	}

	file := parser.Parse()
	return file, parser.Error()
}

func MoFromReader(r io.Reader, name string) (*po.File, error) {
	parser, err := NewMoFromReader(r, name)
	if err != nil {
		return nil, err
	}
	file := parser.Parse()
	return file, parser.Error()
}

func MoFromFile(f *os.File) (*po.File, error) {
	parser, err := NewMoFromFile(f)
	if err != nil {
		return nil, err
	}
	file := parser.Parse()
	return file, parser.Error()
}

func MoFromBytes(b []byte, name string) (*po.File, error) {
	parser := NewMoFromBytes(b, name)
	file := parser.Parse()
	return file, parser.Error()
}
