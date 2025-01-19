package goparse

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strconv"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/poconfig"
	"github.com/Tom5521/xgotext/pkg/poentry"
)

type File struct {
	config    poconfig.Config
	file      *ast.File
	content   []byte
	path      string
	pkgName   string
	hasGotext bool
}

const WantedImport = `"github.com/leonelquinteros/gotext"`

func NewFileFromReader(r io.Reader, name string, cfg poconfig.Config) (*File, error) {
	err := validateConfig(cfg)
	if err != nil {
		return nil, err
	}
	return unsafeNewFileFromReader(r, name, cfg)
}

func unsafeNewFileFromReader(r io.Reader, name string, cfg poconfig.Config) (*File, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return unsafeNewFile(content, name, cfg)
}

func NewFileFromPath(path string, cfg poconfig.Config) (*File, error) {
	err := validateConfig(cfg)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return unsafeNewFileFromReader(file, path, cfg)
}

func NewFileFromBytes(b []byte, name string, cfg poconfig.Config) (*File, error) {
	err := validateConfig(cfg)
	if err != nil {
		return nil, err
	}

	return unsafeNewFile(b, name, cfg)
}

func unsafeNewFile(content []byte, name string, cfg poconfig.Config) (*File, error) {
	file := &File{
		content: content,
		path:    name,
		config:  cfg,
	}
	var err error
	file.file, err = parser.ParseFile(token.NewFileSet(), file.path, content, 0)
	if err != nil {
		return nil, err
	}

	file.determinePackageInfo()
	return file, nil
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
