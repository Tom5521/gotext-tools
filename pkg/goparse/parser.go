package goparse

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	krfs "github.com/kr/fs"

	"github.com/Tom5521/xgotext/pkg/poconfig"
	"github.com/Tom5521/xgotext/pkg/poentry"
)

type Parser struct {
	Config poconfig.Config

	files []*File
	seen  map[string]bool
}

func NewParser(path string, cfg poconfig.Config) (*Parser, error) {
	err := validateConfig(cfg)
	if err != nil {
		return nil, err
	}
	p := baseParser(cfg)
	walker := krfs.Walk(path)
	for walker.Step() {
		if p.shouldSkipFile(walker.Path(), walker.Stat(), walker.Err()) {
			continue
		}

		f, err := unsafeNewFileFromPath(walker.Path(), &cfg)
		if err != nil {
			return nil, err
		}
		p.files = append(p.files, f)
	}

	return p, nil
}

func baseParser(cfg poconfig.Config) *Parser {
	return &Parser{
		Config: cfg,
		seen:   make(map[string]bool),
	}
}

func NewParserFromReader(r io.Reader, name string, cfg poconfig.Config) (*Parser, error) {
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

func unsafeNewParserFromBytes(b []byte, name string, cfg poconfig.Config) (*Parser, error) {
	p := baseParser(cfg)
	f, err := unsafeNewFile(b, name, &cfg)
	if err != nil {
		return nil, err
	}
	p.files = append(p.files, f)

	return p, nil
}

func NewParserFromBytes(b []byte, name string, cfg poconfig.Config) (*Parser, error) {
	err := validateConfig(cfg)
	if err != nil {
		return nil, err
	}
	return unsafeNewParserFromBytes(b, name, cfg)
}

func NewParserFromFile(file *os.File, cfg poconfig.Config) (*Parser, error) {
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

func NewParserFromFiles(files []string, cfg poconfig.Config) (*Parser, error) {
	err := validateConfig(cfg)
	if err != nil {
		return nil, err
	}
	p := baseParser(cfg)
	for _, file := range files {
		walker := krfs.Walk(file)
		for walker.Step() {
			if p.shouldSkipFile(walker.Path(), walker.Stat(), walker.Err()) {
				continue
			}

			f, err := unsafeNewFileFromPath(walker.Path(), &cfg)
			if err != nil {
				return nil, err
			}

			p.files = append(p.files, f)
		}
	}

	return p, nil
}

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
