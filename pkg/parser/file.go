package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strconv"
)

type File struct {
	config    Config
	file      *ast.File
	content   string
	path      string
	pkgName   string
	hasGotext bool
}

const WantedImport = `"github.com/leonelquinteros/gotext"`

func NewFile(path string, cfg Config) (*File, error) {
	cfgErrs := cfg.Validate()
	if len(cfgErrs) > 0 {
		return nil, fmt.Errorf("configuration is invalid: %w", cfgErrs[0])
	}
	return unsafeNewFile(path, cfg)
}

func unsafeNewFile(path string, cfg Config) (*File, error) {
	file := &File{
		path:   path,
		config: cfg,
	}
	if err := file.generateAST(); err != nil {
		return nil, err
	}
	file.determinePackageInfo()
	return file, nil
}

func (f *File) generateAST() error {
	data, err := os.ReadFile(f.path)
	if err != nil {
		return err
	}
	f.content = string(data)

	fset := token.NewFileSet()
	f.file, err = parser.ParseFile(fset, f.path, data, 0)
	return err
}

func (f *File) determinePackageInfo() {
	f.pkgName = "gotext"
	for _, imprt := range f.file.Imports {
		f.hasGotext = imprt.Path.Value == WantedImport
		if f.hasGotext {
			if imprt.Name != nil {
				f.pkgName = imprt.Name.String()
			}
			break
		}
	}
}

func (f *File) shouldSkip(n ast.Node) bool {
	if callExpr, ok := n.(*ast.CallExpr); ok {
		if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			if ident, ok := selectorExpr.X.(*ast.Ident); ok {
				if ident.Name != f.pkgName {
					return true
				}
				if _, ok = gotextGetter[selectorExpr.Sel.Name]; ok {
					return false
				}
			}
		}
	}

	return true
}

func (f *File) extract() (translations []Translation, errs []error) {
	for node := range InspectNode(f.file) {
		if f.shouldSkip(node) {
			if f.config.ExtractAll {
				t, e := f.processString(node)
				translations = append(translations, t...)
				errs = append(errs, e...)
			}
			continue
		}

		callExpr := node.(*ast.CallExpr)
		selectorExpr := callExpr.Fun.(*ast.SelectorExpr)

		translation, valid, err := f.processMethod(selectorExpr.Sel.Name, callExpr)
		if err != nil {
			errs = append(
				errs,
				fmt.Errorf("[%s:%d] %w", f.path, findLine(f.content, selectorExpr.Pos()), err),
			)
		}
		if valid && err == nil {
			translations = append(translations, translation)
		}
	}

	translations = cleanDuplicates(translations)
	return
}

func (f *File) processString(n ast.Node) (translations []Translation, errs []error) {
	if basicLit, ok := n.(*ast.BasicLit); ok {
		if basicLit.Kind == token.STRING {
			line := findLine(f.content, basicLit.Pos())
			content, err := strconv.Unquote(basicLit.Value)
			if err != nil {
				errs = append(
					errs,
					fmt.Errorf(
						"error extracting string in %s line %d: %w",
						f.path,
						line,
						err,
					),
				)
				return
			}
			translations = append(translations,
				Translation{
					ID: content,
					Locations: []Location{
						{line, f.path},
					},
				},
			)
		}
	}

	return
}

func (f *File) ParseTranslations() (translations []Translation, errs []error) {
	return f.extract()
}
