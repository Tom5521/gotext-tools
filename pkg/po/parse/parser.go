package parse

import (
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/Tom5521/xgotext/pkg/po"
	"github.com/antlr4-go/antlr/v4"
)

type Parser struct {
	config  Config
	options []Option

	data     []byte
	filename string

	errors []error
}

func (p *Parser) applyOptions(opts ...Option) {
	for _, opt := range opts {
		opt(&p.config)
	}
}

func NewParser(path string, options ...Option) (*Parser, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return NewParserFromBytes(file, path, options...)
}

func NewParserFromReader(r io.Reader, name string, options ...Option) (*Parser, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return NewParserFromBytes(data, name, options...)
}

func NewParserFromFile(f *os.File, options ...Option) (*Parser, error) {
	return NewParserFromReader(f, f.Name(), options...)
}

func NewParserFromString(s, name string, options ...Option) (*Parser, error) {
	return NewParserFromBytes([]byte(s), name, options...)
}

func NewParserFromBytes(data []byte, name string, options ...Option) (*Parser, error) {
	p := &Parser{
		options:  options,
		config:   DefaultConfig(),
		data:     data,
		filename: name,
	}

	return p, nil
}

func (p Parser) Errors() []error {
	return p.errors
}

func (p *Parser) Parse(options ...Option) *po.File {
	p.applyOptions(p.options...)
	p.applyOptions(options...)

	is := antlr.NewInputStream(string(p.data))

	errListener := &CustomErrorListener{}
	lexer := NewPoLexer(is)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errListener)

	stream := antlr.NewCommonTokenStream(lexer, antlr.LexerDefaultTokenChannel)

	parser := NewPoParser(stream)
	parser.RemoveErrorListeners()
	parser.AddErrorListener(errListener)
	tree := parser.Start_()

	for _, err := range errListener.Errors {
		p.errors = append(p.errors, fmt.Errorf("error in file %s: %s", p.filename, err))
	}

	walker := antlr.NewParseTreeWalker()
	var listener Listener
	walker.Walk(&listener, tree)

	for _, err := range listener.errors {
		err = fmt.Errorf("error in file %s: %w", p.filename, err)
		p.errors = append(p.errors, err)
	}

	for _, err := range p.errors {
		p.config.Logger.Println("ERROR:", err)
	}

	if p.config.SkipHeader {
		i := listener.entries.Index("", "")
		if i != -1 {
			listener.entries = slices.Delete(listener.entries, i, i+1)
		}
	}
	if p.config.CleanDuplicates {
		listener.entries = listener.entries.CleanDuplicates()
	}

	return &po.File{
		Entries: listener.entries,
		Name:    p.filename,
	}
}
