// Package goparse provides tools to parse and process Go source files,
// extracting translations and handling various configurations.
package goparse

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"

	krfs "github.com/kr/fs"

	"github.com/Tom5521/xgotext/pkg/po/config"
	"github.com/Tom5521/xgotext/pkg/po/types"
)

// Parser represents a parser that processes Go files according to a given configuration.
type Parser struct {
	Config config.Config   // Configuration settings for parsing.
	files  []*File         // List of parsed files.
	seen   map[string]bool // Tracks already processed files to avoid duplication.
}

func (p *Parser) appendFiles(files ...string) error {
	for _, file := range files {
		walker := krfs.Walk(file)
		for walker.Step() {
			if p.shouldSkipFile(walker) {
				continue
			}
			if p.Config.Verbose && p.Config.Logger != nil {
				p.Config.Logger.Printf("Reading %s...", walker.Path())
			}
			f, err := NewFileFromPath(walker.Path(), &p.Config)
			if err != nil {
				return fmt.Errorf("error reading file %s: %w", walker.Path(), err)
			}
			p.files = append(p.files, f)
		}
	}

	return nil
}

// NewParser initializes a new Parser for a given directory path and configuration.
func NewParser(path string, cfg config.Config) (*Parser, error) {
	err := validateConfig(cfg)
	if err != nil {
		return nil, err
	}
	p := baseParser(cfg)
	err = p.appendFiles(path)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// baseParser creates a base Parser instance with the provided configuration.
func baseParser(cfg config.Config) *Parser {
	return &Parser{
		Config: cfg,
		seen:   make(map[string]bool),
	}
}

// NewParserFromReader creates a Parser from an io.Reader, such as a file or memory buffer.
func NewParserFromReader(
	r io.Reader,
	name string,
	cfg config.Config,
) (*Parser, error) {
	err := validateConfig(cfg)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return unsafeNewParserFromBytes(data, name, cfg)
}

// unsafeNewParserFromBytes creates a Parser from raw byte data without validating the configuration.
func unsafeNewParserFromBytes(
	b []byte,
	name string,
	cfg config.Config,
) (*Parser, error) {
	p := baseParser(cfg)
	f, err := unsafeNewFile(b, name, &cfg)
	if err != nil {
		return nil, err
	}
	p.files = append(p.files, f)

	return p, nil
}

func NewParserFromString(
	s string,
	name string,
	cfg config.Config,
) (*Parser, error) {
	err := validateConfig(cfg)
	if err != nil {
		return nil, err
	}
	return unsafeNewParserFromBytes([]byte(s), name, cfg)
}

// NewParserFromBytes creates a Parser from raw byte data after validating the configuration.
func NewParserFromBytes(
	b []byte,
	name string,
	cfg config.Config,
) (*Parser, error) {
	err := validateConfig(cfg)
	if err != nil {
		return nil, err
	}
	return unsafeNewParserFromBytes(b, name, cfg)
}

// NewParserFromFile creates a Parser from an os.File instance.
func NewParserFromFile(file *os.File, cfg config.Config) (*Parser, error) {
	err := validateConfig(cfg)
	if err != nil {
		return nil, err
	}
	p := baseParser(cfg)
	f, err := unsafeNewFileFromReader(file, file.Name(), &cfg)
	if err != nil {
		return nil, err
	}

	p.files = append(p.files, f)

	return p, nil
}

// NewParserFromFiles initializes a Parser from a list of file paths.
func NewParserFromFiles(files []string, cfg config.Config) (*Parser, error) {
	err := validateConfig(cfg)
	if err != nil {
		return nil, err
	}
	p := baseParser(cfg)
	err = p.appendFiles(files...)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Parse processes all files associated with the Parser and extracts translations.
func (p *Parser) Parse() (translations []types.Entry, errs []error) {
	for _, f := range p.files {
		t, e := f.Translations()
		if p.Config.Logger != nil {
			for _, err := range e {
				p.Config.Logger.Println(err)
			}
		}
		errs = append(errs, e...)
		translations = append(translations, t...)
	}

	return
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
		if p.Config.Logger != nil {
			p.Config.Logger.Printf("warning: skipping duplicated file: %s", w.Path())
		}
		return true
	}
	p.seen[abs] = true

	return p.isExcludedPath(w.Path())
}

// isExcludedPath checks if a path is in the exclude list defined in the configuration.
func (p Parser) isExcludedPath(path string) bool {
	return slices.ContainsFunc(p.Config.Exclude, func(s string) bool {
		abs1, err1 := filepath.Abs(s)
		abs2, err2 := filepath.Abs(path)
		return (abs1 == abs2) && (err1 == nil && err2 == nil)
	})
}
