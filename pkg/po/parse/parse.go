// Package parse provides the functionality to parse PO (Portable Object)
// files from various sources, such as files, byte slices, strings, and readers.
// It integrates with the `ast`, `generator`, and `types`
// packages to normalize and generate structured representations of PO files.
//
// Key Features:
// - Parses PO files from different input sources (file paths, strings, byte slices, and readers).
// - Normalizes PO entries using the `ast.Normalizer`.
// - Generates a structured `types.File` representation of the PO file.
// - Handles errors and warnings during parsing and normalization.
//
// Example Usage:
//
//	cfg := parsers.Config{}
//	p, err := parse.NewParser("example.po", cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	file := p.Parse()
//	if len(p.Errors()) > 0 {
//	    log.Fatal(p.Errors()[0])
//	}
//	fmt.Println(file)
//
// For more details, refer to the individual functions and types.
package parse

import (
	"io"
	"os"

	"github.com/Tom5521/xgotext/pkg/po"
	"github.com/Tom5521/xgotext/pkg/po/parse/ast"
	"github.com/Tom5521/xgotext/pkg/po/parse/generator"
)

type Parser struct {
	config  Config
	options []Option
	norm    *ast.Normalizer

	warns  []string
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
		options: options,
		config:  DefaultConfig(),
	}
	err := p.processpath(data, name)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p Parser) Errors() []error {
	return p.errors
}

func (p Parser) Warnings() []string {
	return p.warns
}

func (p *Parser) processpath(content []byte, path string) error {
	if p.config.Verbose {
		p.config.Logger.Println("Extracting tokens...")
	}
	norm, errs := ast.NewParser(content, path).Normalizer()
	for _, e := range errs {
		p.config.Logger.Println("ERROR:", e)
	}
	if len(errs) > 0 {
		return errs[0]
	}

	p.norm = norm
	return nil
}

func (p *Parser) Parse(options ...Option) *po.File {
	p.applyOptions(p.options...)
	p.applyOptions(options...)
	p.errors = nil // Reset errors.
	p.warns = nil  // Reset warnings.

	if p.config.Verbose {
		p.config.Logger.Println("Parsing...")
	}

	p.norm.Normalize()
	for _, w := range p.norm.Warnings() {
		p.config.Logger.Println("WARN:", w)
	}
	for _, e := range p.norm.Errors() {
		p.config.Logger.Println("ERROR:", e)
	}
	p.warns = append(p.warns, p.norm.Warnings()...)
	if len(p.norm.Errors()) > 0 {
		p.errors = append(p.errors, p.norm.Errors()...)
		return nil
	}

	g := generator.New(p.norm.File())

	file := g.Generate()
	if len(g.Errors()) > 0 {
		p.errors = append(p.errors, g.Errors()...)
		return nil
	}

	return file
}
