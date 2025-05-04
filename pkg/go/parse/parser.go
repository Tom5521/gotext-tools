// Package goparse provides tools to parse and process Go source files,
// extracting translations and handling various configurations.
package parse

import (
	"errors"
	"fmt"
	"io"
	"os"

	krfs "github.com/kr/fs"

	"github.com/Tom5521/gotext-tools/v2/internal/util"
	"github.com/Tom5521/gotext-tools/v2/pkg/po"
)

var _ po.Parser = (*Parser)(nil)

type Parser struct {
	Config Config              // Configuration settings for parsing.
	files  []*File             // List of parsed files.
	seen   map[string]struct{} // Tracks already processed files to avoid duplication.

	errors []error
}

func (p *Parser) error(format string, a ...any) {
	var err error
	format = "go/parse: " + format
	if len(a) == 0 {
		err = errors.New(format)
	} else {
		err = fmt.Errorf(format, a...)
	}

	if p.Config.Logger != nil {
		p.Config.Logger.Println("ERROR:", err)
	}

	p.errors = append(p.errors, err)
}

func (p *Parser) lastErr() error {
	if len(p.errors) == 0 {
		return nil
	}

	return p.errors[len(p.errors)-1]
}

func (p *Parser) info(format string, a ...any) {
	if p.Config.Logger != nil && p.Config.Verbose {
		p.Config.Logger.Println("INFO:", fmt.Sprintf(format, a...))
	}
}

func (p *Parser) appendFiles(files ...string) error {
	for _, file := range files {
		walker := krfs.Walk(file)
		for walker.Step() {
			if util.ShouldSkipFile(walker, p.Config.Exclude, &p.seen, p.Config.Logger) {
				continue
			}

			p.info("Reading %s...", walker.Path())
			f, err := NewFileFromPath(walker.Path(), &p.Config)
			if err != nil {
				p.error("error reading file %s: %w", walker.Path(), err)
				return p.lastErr()
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
		p.error("error parsing files: %w", err)
		return nil, p.lastErr()
	}

	return p, nil
}

// baseParser creates a base Parser instance with the provided configuration.
func baseParser(options ...Option) *Parser {
	p := &Parser{
		Config: DefaultConfig(options...),
		seen:   make(map[string]struct{}),
	}

	return p
}

// NewParserFromReader creates a Parser from an io.Reader, such as a file or memory buffer.
func NewParserFromReader(
	r io.Reader,
	name string,
	options ...Option,
) (*Parser, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		p := &Parser{Config: DefaultConfig(options...)}
		p.error("error reading: %w", err)
		return nil, p.lastErr()
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
	f, err := NewFileFromBytes(b, name, &p.Config)
	if err != nil {
		p.error("error configuring file: %w", err)
		return nil, p.lastErr()
	}
	p.files = append(p.files, f)

	return p, nil
}

func NewParserFromFiles(files []*os.File, options ...Option) (*Parser, error) {
	p := baseParser(options...)
	for _, file := range files {
		f, err := NewFile(file, file.Name(), &p.Config)
		if err != nil {
			p.error("error configuring file: %w", err)
			return nil, p.lastErr()
		}
		f.config = &p.Config
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
		p.error("error parsing files: %w", err)
		return nil, p.lastErr()
	}

	return p, nil
}

func (p *Parser) ParseWithOptions(options ...Option) *po.File {
	p.Config.ApplyOptions(options...)
	defer p.Config.RestoreLastCfg()
	return p.Parse()
}

// Parse processes all files associated with the Parser and extracts translations.
func (p *Parser) Parse() (file *po.File) {
	file = new(po.File)
	p.errors = nil // Clean errors

	if !p.Config.NoHeader {
		header := po.DefaultTemplateHeader()
		if p.Config.Header != nil {
			header = *p.Config.Header
		}

		if p.Config.HeaderConfig != nil {
			header = p.Config.HeaderConfig.ToHeaderWithDefaults()
		}

		if p.Config.HeaderOptions != nil {
			header = po.HeaderConfigFromOptions(p.Config.HeaderOptions...).ToHeaderWithDefaults()
		}

		header.Fields = append(header.Fields, po.HeaderField{Key: "X-Generator", Value: "xgotext"})

		file.Entries = append(file.Entries, header.ToEntry())
	}

	for _, f := range p.files {
		p.info("parsing %s...", f.name)
		entries := f.Entries()
		if err := f.Error(); err != nil {
			continue
		}
		file.Entries = append(file.Entries, entries...)
	}

	if p.Config.CleanDuplicates {
		file.Entries = file.CleanDuplicates()
	}

	return
}

func (p Parser) Error() error {
	if len(p.errors) == 0 {
		return nil
	}

	return p.errors[0]
}

func (p Parser) Errors() []error {
	return p.errors
}

// Files returns the list of files associated with the Parser.
func (p Parser) Files() []*File {
	return p.files
}
