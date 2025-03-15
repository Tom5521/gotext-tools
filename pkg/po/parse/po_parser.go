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
	Config PoConfig

	data     string
	filename string

	errors []error
}

func NewPo(path string, options ...PoOption) (*PoParser, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return NewPoFromBytes(file, path, options...), nil
}

func NewPoFromReader(r io.Reader, name string, options ...PoOption) (*PoParser, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return NewPoFromBytes(data, name, options...), nil
}

func NewPoFromFile(f *os.File, options ...PoOption) (*PoParser, error) {
	return NewPoFromReader(f, f.Name(), options...)
}

func NewPoFromString(s, name string, options ...PoOption) *PoParser {
	return &PoParser{
		Config:   DefaultPoConfig(options...),
		data:     s,
		filename: name,
	}
}

func NewPoFromBytes(data []byte, name string, options ...PoOption) *PoParser {
	return NewPoFromString(string(data), name, options...)
}

// Return the first error in the stack.
func (p PoParser) Error() error {
	if len(p.errors) == 0 {
		return nil
	}
	return p.errors[0]
}

func (p PoParser) Errors() []error {
	return p.errors
}

func (p *PoParser) Parse(options ...PoOption) *po.File {
	p.Config.ApplyOptions(options...)
	defer p.Config.RestoreLastCfg()

	p.errors = nil

	is := antlr.NewInputStream(p.data)

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
		p.Config.Logger.Println("ERROR:", err)
	}

	if p.Config.SkipHeader {
		i := listener.Entries.Index("", "")
		if i != -1 {
			listener.Entries = slices.Delete(listener.Entries, i, i+1)
		}
	}
	if p.Config.CleanDuplicates {
		listener.Entries = listener.Entries.CleanDuplicates()
	}

	return &po.File{
		Entries: listener.Entries,
		Name:    p.filename,
	}
}
