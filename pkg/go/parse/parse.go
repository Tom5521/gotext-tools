package parse

import (
	"io"
	"os"

	"github.com/Tom5521/gotext-tools/pkg/po"
)

func FromPath(path string, opts ...Option) (*po.File, error) {
	parser, err := NewParser(path, opts...)
	if err != nil {
		return nil, err
	}
	file := parser.Parse()
	return file, parser.Error()
}

func FromPaths(paths []string, opts ...Option) (*po.File, error) {
	parser, err := NewParserFromPaths(paths, opts...)
	if err != nil {
		return nil, err
	}
	file := parser.Parse()
	return file, parser.Error()
}

func FromReader(r io.Reader, name string, opts ...Option) (*po.File, error) {
	parser, err := NewParserFromReader(r, name, opts...)
	if err != nil {
		return nil, err
	}
	file := parser.Parse()
	return file, parser.Error()
}

func FromString(s, name string, opts ...Option) (*po.File, error) {
	parser, err := NewParserFromString(s, name, opts...)
	if err != nil {
		return nil, err
	}
	file := parser.Parse()
	return file, parser.Error()
}

func FromBytes(b []byte, name string, opts ...Option) (*po.File, error) {
	parser, err := NewParserFromBytes(b, name, opts...)
	if err != nil {
		return nil, err
	}
	file := parser.Parse()
	return file, parser.Error()
}

func FromFiles(files []*os.File, opts ...Option) (*po.File, error) {
	parser, err := NewParserFromFiles(files, opts...)
	if err != nil {
		return nil, err
	}
	file := parser.Parse()
	return file, parser.Error()
}

func FromFile(f *os.File, opts ...Option) (*po.File, error) {
	parser, err := NewParserFromFile(f, opts...)
	if err != nil {
		return nil, err
	}
	file := parser.Parse()
	return file, parser.Error()
}
