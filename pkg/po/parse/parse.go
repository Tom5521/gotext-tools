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
	p := baseParser(cfg)
	p.file, err = processpath(data, name)
	return p, err
}

func processpath(content []byte, path string) (*ast.File, error) {
	parser := ast.NewParser(content, path)
	errs := parser.Parse()
	if len(errs) > 0 {
		return nil, errs[0]
	}

	return parser.File, nil
}

func (p *Parser) Parse() (*types.File, []string, []error) {
	g := generator.New(p.file)

	if len(g.Warnings()) > 0 && p.Config.Logger != nil {
		for _, warn := range g.Warnings() {
			p.Config.Logger.Print("WARN: ", warn)
		}
	}

	return g.Generate(), g.Warnings(), g.Errors()
}
