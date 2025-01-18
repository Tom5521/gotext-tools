package parser

import (
	_ "embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"
)

//go:embed template.pot
var PotHeader string

type Parser struct {
	Config Config

	translations []Translation
	files        []*File
	seen         map[string]bool
}

func NewParser(cfg Config) (p *Parser, err error) {
	cfgErrs := cfg.Validate()
	if len(cfgErrs) > 0 {
		return nil, fmt.Errorf("configuration is invalid: %w", cfgErrs[0])
	}
	p = &Parser{
		Config: cfg,
		seen:   make(map[string]bool),
	}
	for _, path := range cfg.Files {
		err = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
			if p.shouldSkipFile(path, info, err) {
				return nil
			}

			file, err := unsafeNewFile(path, cfg)
			if err != nil {
				return err
			}
			p.files = append(p.files, file)

			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return
}

func (p *Parser) Parse() (errs []error) {
	for _, f := range p.files {
		t, e := f.ParseTranslations()
		errs = append(errs, e...)
		p.translations = append(p.translations, t...)
	}

	return
}

func (p Parser) Files() []*File {
	return p.files
}

func (p Parser) Compile() []byte {
	var b strings.Builder
	fmt.Fprintf(&b, PotHeader, p.Config.PackageVersion, p.Config.Language, p.Config.Nplurals)
	var translations []Translation
	copy(translations, p.translations)

	translations = cleanDuplicates(translations)
	for _, t := range translations {
		fmt.Fprintln(&b, t.Format(p.Config.Nplurals))
	}
	return []byte(b.String())
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
