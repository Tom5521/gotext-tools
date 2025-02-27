// Package goparse provides tools to parse and process Go source files,
// extracting translations and handling various configurations.
package parse

import (
	"fmt"
	"io"
	"os"

	krfs "github.com/kr/fs"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po"
)

// Parser represents a parser for processing Go source files and extracting translation entries.
//
// ### Attributes:
// - `config`: Configuration settings for the parser (e.g., exclude paths, verbose logging).
// - `options`: Additional options to customize the parser behavior.
// - `files`: A list of `File` objects representing the Go source files to be processed.
// - `seen`: Tracks already processed files to avoid duplicate processing.
// - `errors`: Stores errors encountered during parsing.
//
// ### Responsibilities:
// - Manage the parsing process, including file handling, configuration, and error reporting.
// - Traverse ASTs to extract translation entries and generate compatible PO file data.
//
// ### Methods:
// - `NewParser`: Initializes a parser for a directory of Go files.
// - `NewParserFromReader`: Creates a parser from an `io.Reader` (e.g., file or memory buffer).
// - `NewParserFromString`: Creates a parser from a string containing Go source code.
// - `NewParserFromBytes`: Creates a parser from raw byte data.
// - `NewParserFromFiles`: Initializes a parser for a list of file paths.
// - `Parse`: Processes all files in the parser, extracting translations and generating entries.
// - `Errors`: Returns any errors encountered during parsing.
// - `Files`: Returns the list of files associated with the parser.
type Parser struct {
	config  Config // Configuration settings for parsing.
	options []Option
	files   []*File         // List of parsed files.
	seen    map[string]bool // Tracks already processed files to avoid duplication.

	errors []error
}

func (p *Parser) applyOptions(options ...Option) {
	for _, opt := range options {
		opt(&p.config)
	}
}

func (p *Parser) appendFiles(files ...string) error {
	for _, file := range files {
		walker := krfs.Walk(file)
		for walker.Step() {
			if util.ShouldSkipFile(walker, p.config.Exclude, &p.seen, p.config.Logger) {
				continue
			}

			if p.config.Verbose {
				p.config.Logger.Println("Reading", walker.Path(), "...")
			}
			f, err := NewFileFromPath(walker.Path(), p.options...)
			if err != nil {
				err = fmt.Errorf("error reading file %s: %w", walker.Path(), err)
				p.config.Logger.Println("ERROR:", err.Error())
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
		err = fmt.Errorf("error parsing files: %w", err)
		p.config.Logger.Println("ERROR:", err)
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
		logger.Println("ERROR:", err)
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
	f, err := NewFile(b, name, options...)
	if err != nil {
		err = fmt.Errorf("error configuring file: %w", err)
		p.config.Logger.Println("ERROR:", err)
		return nil, err
	}
	p.files = append(p.files, f)

	return p, nil
}

func NewParserFromFiles(files []*os.File, options ...Option) (*Parser, error) {
	p := baseParser(options...)
	for _, file := range files {
		f, err := NewFileFromReader(file, file.Name(), options...)
		if err != nil {
			err = fmt.Errorf("error configuring file: %w", err)
			p.config.Logger.Println("ERROR:", err)
			return nil, err
		}
		p.files = append(p.files, f)
	}

	return p, nil
}

// NewParserFromFile creates a Parser from an os.File instance.
func NewParserFromFile(file *os.File, options ...Option) (*Parser, error) {
	return NewParserFromFiles([]*os.File{file}, options...)
}

// NewParserFromPaths initializes a Parser from a list of file paths.
func NewParserFromPaths(files []string, options ...Option) (*Parser, error) {
	p := baseParser(options...)
	err := p.appendFiles(files...)
	if err != nil {
		err = fmt.Errorf("error parsing files: %w", err)
		p.config.Logger.Println("ERROR:", err)
		return nil, err
	}

	return p, nil
}

// Parse processes all files associated with the Parser and extracts translations.
func (p *Parser) Parse(options ...Option) (file *po.File) {
	p.applyOptions(p.options...)
	p.applyOptions(options...)
	defer p.applyOptions(p.options...) // Reset default settings.
	file = &po.File{}
	p.errors = nil // Clean errors

	header := po.DefaultHeader()
	if p.config.Header != nil {
		header = *p.config.Header
	}

	if p.config.HeaderConfig != nil {
		header = p.config.HeaderConfig.ToHeaderWithDefaults()
	}

	if p.config.HeaderOptions != nil {
		header = po.HeaderConfigFromOptions(p.config.HeaderOptions...).ToHeaderWithDefaults()
	}

	header.Fields = append(header.Fields, po.HeaderField{Key: "X-Generator", Value: "xgotext"})

	file.Entries = append(file.Entries, header.ToEntry())

	for _, f := range p.files {
		if p.config.Verbose {
			p.config.Logger.Println("Parsing", f.path, "...")
		}
		entries, e := f.Entries()
		if len(e) > 0 {
			p.errors = append(p.errors, e...)
			for _, err := range e {
				p.config.Logger.Println(
					"ERROR:",
					fmt.Errorf("error parsing file %s: %w", f.path, err),
				)
			}
			continue
		}
		file.Entries = append(file.Entries, entries...)
	}

	if p.config.CleanDuplicates {
		file.Entries = file.Entries.CleanDuplicates()
	}

	return
}

func (p Parser) Errors() []error {
	return p.errors
}

// Files returns the list of files associated with the Parser.
func (p Parser) Files() []*File {
	return p.files
}
