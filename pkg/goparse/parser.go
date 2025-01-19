package goparse

import (
	_ "embed"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"slices"

	"github.com/Tom5521/xgotext/pkg/poconfig"
	"github.com/Tom5521/xgotext/pkg/poentry"
)

type Parser struct {
	Config poconfig.Config

	files []*File
	seen  map[string]bool
}

func NewParser(cfg poconfig.Config) (p *Parser, err error) {
	return
}

func newBaseParser(cfg poconfig.Config) (*Parser, error) {
	cfgErrs := cfg.Validate()
	if len(cfgErrs) > 0 {
		return nil, fmt.Errorf("configuration is invalid: %w", cfgErrs[0])
	}
	return &Parser{
		Config: cfg,
		seen:   make(map[string]bool),
	}, nil
}

func NewParserFromReader(r io.Reader, cfg poconfig.Config) (*Parser, error)
func NewParserFromBytes(b []byte, cfg poconfig.Config) (*Parser, error)
func NewParserFromFiles(files []string, cfg poconfig.Config) (*Parser, error)

func (p *Parser) Parse() (translations []poentry.Translation, errs []error) {
	for _, f := range p.files {
		t, e := f.ParseTranslations()
		errs = append(errs, e...)
		translations = append(translations, t...)
	}

	return
}

func (p Parser) Files() []*File {
	return p.files
}

// shouldSkipFile determines if a file should be skipped during processing.
func (p Parser) shouldSkipFile(path string, info fs.FileInfo, err error) bool {
	if err != nil || info.IsDir() {
		return true
	}

	if filepath.Ext(path) != ".go" {
		return true
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return true
	}

	_, seen := p.seen[abs]
	if seen {
		return true
	}
	p.seen[abs] = true

	return p.isExcludedPath(path)
}

// isExcludedPath checks if a path is in the exclude list.
func (p Parser) isExcludedPath(path string) bool {
	return slices.ContainsFunc(p.Config.Exclude, func(s string) bool {
		abs1, err1 := filepath.Abs(s)
		abs2, err2 := filepath.Abs(path)
		return (abs1 == abs2) && (err1 == nil && err2 == nil)
	})
}
