package parse

import (
	"io"
	"os"

	"github.com/Tom5521/xgotext/pkg/po/config"
	"github.com/Tom5521/xgotext/pkg/po/parse/ast"
	"github.com/Tom5521/xgotext/pkg/po/parse/generator"
	"github.com/Tom5521/xgotext/pkg/po/types"
)

type Parser struct {
	Config config.Config
	file   *ast.File
	seen   map[string]bool
	warns  []string
}

func baseParser(cfg config.Config) *Parser {
	return &Parser{
		Config: cfg,
		seen:   make(map[string]bool),
	}
}

func NewParserFromReader(r io.Reader, name string, cfg config.Config) (*Parser, error) {
	err := validateConfig(cfg)
	if err != nil {
		return nil, err
	}
	return unsafeNewParserFromReader(r, name, cfg)
}

func unsafeNewParserFromReader(r io.Reader, name string, cfg config.Config) (*Parser, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return unsafeNewParserFromBytes(data, name, cfg)
}

func NewParserFromFile(f *os.File, cfg config.Config) (*Parser, error) {
	err := validateConfig(cfg)
	if err != nil {
		return nil, err
	}

	return unsafeNewParserFromReader(f, f.Name(), cfg)
}

func NewParserFromString(s, name string, cfg config.Config) (*Parser, error) {
	err := validateConfig(cfg)
	if err != nil {
		return nil, err
	}
	return unsafeNewParserFromBytes([]byte(s), name, cfg)
}

func NewParserFromBytes(d []byte, name string, cfg config.Config) (*Parser, error) {
	err := validateConfig(cfg)
	if err != nil {
		return nil, err
	}
	return unsafeNewParserFromBytes(d, name, cfg)
}

func unsafeNewParserFromBytes(data []byte, name string, cfg config.Config) (*Parser, error) {
	var err error
	var warns []string
	p := baseParser(cfg)
	p.file, warns, err = p.processpath(data, name)
	p.warns = append(p.warns, warns...)
	return p, err
}

func (p *Parser) processpath(content []byte, path string) (*ast.File, []string, error) {
	norm, errs := ast.NewParser(content, path).Normalizer()
	if len(errs) > 0 {
		return nil, nil, errs[0]
	}

	norm.Normalize()

	if len(norm.Errors()) > 0 {
		return nil, norm.Warnings(), norm.Errors()[0]
	}

	if p.Config.Logger != nil && p.Config.Verbose {
		for _, warn := range norm.Warnings() {
			p.Config.Logger.Println("WARN:", warn)
		}
	}

	return norm.File(), norm.Warnings(), nil
}

func (p *Parser) Parse() (*types.File, []string, []error) {
	g := generator.New(p.file)

	return g.Generate(), p.warns, g.Errors()
}
