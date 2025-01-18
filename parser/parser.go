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
	OnWalk func(string, fs.FileInfo, error)

	files   []*File
	exclude []string
	seen    map[string]bool
}

func NewParser(exclude []string, paths ...string) (p *Parser, err error) {
	p = &Parser{
		exclude: exclude,
		seen:    make(map[string]bool),
	}
	for _, path := range paths {
		err = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
			if p.shouldSkipFile(path, info, err) {
				return nil
			}

			if p.OnWalk != nil {
				p.OnWalk(path, info, err)
			}

			file, err := NewFile(path)
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
		errs = append(errs, f.ParseTranslations()...)
	}

	return
}

func (p Parser) Files() []*File {
	return p.files
}

func (p Parser) Compile(version, lang string, plurals uint) []byte {
	var b strings.Builder
	fmt.Fprintf(&b, PotHeader, version, lang, plurals)

	var translations []Translation

	for _, f := range p.files {
		translations = append(translations, f.Translations...)
	}
	translations = cleanDuplicates(translations)
	for _, t := range translations {
		fmt.Fprintln(&b, t.Format(plurals))
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
	return slices.ContainsFunc(p.exclude, func(s string) bool {
		abs1, err1 := filepath.Abs(s)
		abs2, err2 := filepath.Abs(path)
		return (abs1 == abs2) && (err1 == nil && err2 == nil)
	})
}
