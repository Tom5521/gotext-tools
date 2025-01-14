package parser

import (
	"bytes"
	_ "embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"

	"github.com/Tom5521/xgotext/flags"
	"github.com/gookit/color"
)

//go:embed template.pot
var PotHeader string

type Parser struct {
	files []*File
}

func NewParser(path string) (p *Parser, err error) {
	p = &Parser{}
	err = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if shouldSkipFile(path, info, err, flags.Exclude) {
			return nil
		}

		if flags.Verbose {
			fmt.Println(path)
		}

		file, err := NewFile(path)
		if err != nil {
			return err
		}
		p.files = append(p.files, file)

		return nil
	})

	return
}

func (p *Parser) Parse() {
	for _, f := range p.files {
		f.ParseTranslations()
	}
}

func (p *Parser) Compile() []byte {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf(PotHeader, flags.ProjectVersion, flags.Language, flags.Nplurals))

	for _, f := range p.files {
		for _, t := range f.Translations {
			fmt.Fprintln(&b, t)
		}
	}
	return b.Bytes()
}

// shouldSkipFile determines if a file should be skipped during processing.
func shouldSkipFile(path string, info fs.FileInfo, err error, exclude []string) bool {
	if err != nil || info.IsDir() {
		return true
	}

	if filepath.Ext(path) != ".go" {
		return true
	}

	return isExcludedPath(path, exclude)
}

// isExcludedPath checks if a path is in the exclude list.
func isExcludedPath(path string, exclude []string) bool {
	return slices.ContainsFunc(exclude, func(s string) bool {
		abs1, err := filepath.Abs(s)
		if err != nil {
			color.Errorln(err)
			return false
		}

		abs2, err := filepath.Abs(path)
		if err != nil {
			color.Errorln(err)
			return false
		}

		return abs1 == abs2
	})
}
