package parse

import (
	"fmt"
	"io"
	"os"
	"slices"

	parser "github.com/Tom5521/xgotext/internal/antlr-po"
	"github.com/Tom5521/xgotext/pkg/po"
	"github.com/antlr4-go/antlr/v4"
)

type PoParser struct {
	config  PoConfig
	options []PoOption

	data     []byte
	filename string

	errors []error
}

func (p *PoParser) applyOptions(opts ...PoOption) {
	for _, opt := range opts {
		opt(&p.config)
	}
}

func NewPo(path string, options ...PoOption) (*PoParser, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return NewPoFromBytes(file, path, options...)
}

func NewPoFromReader(r io.Reader, name string, options ...PoOption) (*PoParser, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return NewPoFromBytes(data, name, options...)
}

func NewPoFromFile(f *os.File, options ...PoOption) (*PoParser, error) {
	return NewPoFromReader(f, f.Name(), options...)
}

func NewPoFromString(s, name string, options ...PoOption) (*PoParser, error) {
	return NewPoFromBytes([]byte(s), name, options...)
}

func NewPoFromBytes(data []byte, name string, options ...PoOption) (*PoParser, error) {
	p := &PoParser{
		options:  options,
		config:   DefaultPoConfig(),
		data:     data,
		filename: name,
	}

	return p, nil
}

func (p PoParser) Errors() []error {
	return p.errors
}

func (p *PoParser) Parse(options ...PoOption) *po.File {
	p.applyOptions(p.options...)
	p.applyOptions(options...)

	p.errors = nil

	is := antlr.NewInputStream(string(p.data))

	errListener := &parser.CustomErrorListener{}
	lexer := parser.NewPoLexer(is)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errListener)

	stream := antlr.NewCommonTokenStream(lexer, antlr.LexerDefaultTokenChannel)

	poParser := parser.NewPoParser(stream)
	poParser.RemoveErrorListeners()
	poParser.AddErrorListener(errListener)
	tree := poParser.Start_()

	for _, err := range errListener.Errors {
		p.errors = append(p.errors, fmt.Errorf("error in file %s: %s", p.filename, err))
	}

	walker := antlr.NewParseTreeWalker()
	var listener parser.Listener
	walker.Walk(&listener, tree)

	for _, err := range listener.Errors {
		err = fmt.Errorf("error in file %s: %w", p.filename, err)
		p.errors = append(p.errors, err)
	}

	for _, err := range p.errors {
		p.config.Logger.Println("ERROR:", err)
	}

	if p.config.SkipHeader {
		i := listener.Entries.Index("", "")
		if i != -1 {
			listener.Entries = slices.Delete(listener.Entries, i, i+1)
		}
	}
	if p.config.CleanDuplicates {
		listener.Entries = listener.Entries.CleanDuplicates()
	}

	return &po.File{
		Entries: listener.Entries,
		Name:    p.filename,
	}
}
