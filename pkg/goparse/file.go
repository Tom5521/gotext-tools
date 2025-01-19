package goparse

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strconv"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/poconfig"
	"github.com/Tom5521/xgotext/pkg/poentry"
)

type File struct {
	config    poconfig.Config
	file      *ast.File
	content   string
	path      string
	pkgName   string
	hasGotext bool
}

const WantedImport = `"github.com/leonelquinteros/gotext"`

func NewFile(path string, cfg poconfig.Config) (*File, error) {
	cfgErrs := cfg.Validate()
	if len(cfgErrs) > 0 {
		return nil, fmt.Errorf("configuration is invalid: %w", cfgErrs[0])
	}
	return unsafeNewFile(path, cfg)
}

func unsafeNewFile(path string, cfg poconfig.Config) (*File, error) {
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

func (f *File) extract() (ts []poentry.Translation, errs []error) {
	for node := range InspectNode(f.file) {
		if f.shouldSkip(node) {
			// This is probably broken.
			// TODO: Fix this.
			if f.config.ExtractAll {
				t, e := f.processString(node)
				ts = append(ts, t...)
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
				fmt.Errorf("[%s:%d] %w", f.path, util.FindLine(f.content, selectorExpr.Pos()), err),
			)
		}
		if valid && err == nil {
			ts = append(ts, translation)
		}
	}

	ts = util.CleanDuplicates(ts)
	return
}

func (f *File) processString(n ast.Node) (ts []poentry.Translation, errs []error) {
	if basicLit, ok := n.(*ast.BasicLit); ok {
		if basicLit.Kind == token.STRING {
			line := util.FindLine(f.content, basicLit.Pos())
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
			ts = append(ts,
				poentry.Translation{
					ID: content,
					Locations: []poentry.Location{
						{line, f.path},
					},
				},
			)
		}
	}

	return
}

func (f *File) ParseTranslations() (translations []poentry.Translation, errs []error) {
	return f.extract()
}
