// Package goparse provides tools to parse and process Go source files,
// extracting translations and handling various configurations.
package parse

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"

	krfs "github.com/kr/fs"

	"github.com/Tom5521/xgotext/pkg/po/types"
)

// Parser represents a parser that processes Go files according to a given configuration.
type Parser struct {
	config  Config // Configuration settings for parsing.
	options []Option
	files   []*File         // List of parsed files.
	seen    map[string]bool // Tracks already processed files to avoid duplication.

	errors []error
}

func (p *Parser) appendFiles(files ...string) error {
	for _, file := range files {
		walker := krfs.Walk(file)
		for walker.Step() {
			if p.shouldSkipFile(walker) {
				continue
			}
			f, err := NewFileFromPath(walker.Path(), p.options...)
			if err != nil {
				err = fmt.Errorf("error reading file %s: %w", walker.Path(), err)
				p.config.Logger.Println(err.Error())
				return err
			}
			p.files = append(p.files, f)
		}
	}

	return nil
}

// NewParser initializes a new Parser for a given directory path and configuration.
func NewParser(path string, options ...Option) (*Parser, error) {
	p := baseParser(options...)
	err := p.appendFiles(path)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// baseParser creates a base Parser instance with the provided configuration.
func baseParser(options ...Option) *Parser {
	p := &Parser{
		options: options,
		config:  DefaultConfig(options...),
		seen:    make(map[string]bool),
	}

	return p
}

// NewParserFromReader creates a Parser from an io.Reader, such as a file or memory buffer.
func NewParserFromReader(
	r io.Reader,
	name string,
	options ...Option,
) (*Parser, error) {
	logger := DefaultConfig(options...).Logger
	data, err := io.ReadAll(r)
	if err != nil {
		err = fmt.Errorf("error reading: %w", err)
		logger.Println(err)
		return nil, err
	}
	return NewParserFromBytes(data, name, options...)
}

func NewParserFromString(
	s string,
	name string,
	options ...Option,
) (*Parser, error) {
	return NewParserFromBytes([]byte(s), name, options...)
}

// NewParserFromBytes creates a Parser from raw byte data after validating the configuration.
func NewParserFromBytes(
	b []byte,
	name string,
	options ...Option,
) (*Parser, error) {
	p := baseParser(options...)
	f, err := NewFileFromBytes(b, name, options...)
	if err != nil {
		err = fmt.Errorf("error configuring file: %w", err)
		p.config.Logger.Println(err)
		return nil, err
	}
	p.files = append(p.files, f)

	return p, nil
}

// NewParserFromFile creates a Parser from an os.File instance.
func NewParserFromFile(file *os.File, options ...Option) (*Parser, error) {
	p := baseParser(options...)
	f, err := NewFileFromReader(file, file.Name(), options...)
	if err != nil {
		err = fmt.Errorf("error configuring file: %w", err)
		p.config.Logger.Println(err.Error())
		return nil, err
	}

	p.files = append(p.files, f)

	return p, nil
}

// NewParserFromFiles initializes a Parser from a list of file paths.
func NewParserFromFiles(files []string, options ...Option) (*Parser, error) {
	p := baseParser(options...)
	err := p.appendFiles(files...)
	if err != nil {
		err = fmt.Errorf("error parsing files: %w", err)
		p.config.Logger.Println(err)
		return nil, err
	}

	return p, nil
}

// Parse processes all files associated with the Parser and extracts translations.
func (p *Parser) Parse() (file *types.File) {
	file = &types.File{}

	header := *p.config.Header

	if p.config.HeaderConfig != nil {
		header = p.config.HeaderConfig.ToHeader()
	}

	if p.config.HeaderOptions != nil {
		header = types.HeaderConfigFromOptions(p.config.HeaderOptions...).ToHeaderWithDefaults()
	}

	file.Entries = append(file.Entries, header.ToEntry())

	for _, f := range p.files {
		entries, e := f.Entries()
		if len(e) > 0 {
			p.errors = append(p.errors, e...)
			for _, err := range e {
				p.config.Logger.Println(fmt.Errorf("error parsing entries: %w", err))
			}
			continue
		}
		file.Entries = append(file.Entries, entries...)
	}

	file.Entries = file.Entries.Solve(p.config.FuzzyMatch)

	return
}

func (p Parser) Errors() []error {
	return p.errors
}

// Files returns the list of files associated with the Parser.
func (p Parser) Files() []*File {
	return p.files
}

// shouldSkipFile determines if a file should be skipped during processing.
func (p *Parser) shouldSkipFile(w *krfs.Walker) bool {
	if w.Err() != nil || w.Stat().IsDir() {
		return true
	}

	if filepath.Ext(w.Path()) != ".go" {
		return true
	}

	abs, err := filepath.Abs(w.Path())
	if err != nil {
		return true
	}

	_, seen := p.seen[abs]
	if seen {
		p.config.Logger.Printf("skipping duplicated file: %s", w.Path())
		return true
	}
	p.seen[abs] = true

	return p.isExcludedPath(w.Path())
}

// isExcludedPath checks if a path is in the exclude list defined in the configuration.
func (p Parser) isExcludedPath(path string) bool {
	return slices.ContainsFunc(p.config.Exclude, func(s string) bool {
		abs1, err1 := filepath.Abs(s)
		abs2, err2 := filepath.Abs(path)
		return (abs1 == abs2) && (err1 == nil && err2 == nil)
	})
}
