package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"slices"
)

type File struct {
	Translations []Translation

	file      *ast.File
	content   string
	path      string
	pkgName   string
	hasGotext bool
}

const WantedImport = `"github.com/leonelquinteros/gotext"`

const (
	Get    = "Get"
	GetD   = "GetD"
	GetN   = "GetN"
	GetC   = "GetC"
	GetND  = "GetND"
	GetNC  = "GetNC"
	GetNDC = "GetNDC"
)

var KnownMethods = []string{
	Get, GetD, GetN, GetC, GetND, GetNC, GetNDC,
}

func NewFile(path string) (*File, error) {
	file := &File{path: path}
	if err := file.generateAST(); err != nil {
		return nil, err
	}
	file.determineHasGotext()
	file.determinePackageName()
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

func (f *File) determineHasGotext() {
	for _, imprt := range f.file.Imports {
		if imprt.Path.Value == WantedImport {
			f.hasGotext = true
			break
		}
	}
}

func (f *File) determinePackageName() {
	f.pkgName = "gotext"
	for _, imprt := range f.file.Imports {
		if imprt.Path.Value == WantedImport {
			if imprt.Name != nil {
				f.pkgName = imprt.Name.String()
			}
			break
		}
	}
}

func (f *File) ParseTranslations() {
	f.Translations = nil
	ast.Inspect(f.file, func(node ast.Node) bool {
		callExpr, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}

		selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		ident, ok := selectorExpr.X.(*ast.Ident)
		if !ok || ident.Name != f.pkgName {
			return true
		}

		if !slices.Contains(KnownMethods, selectorExpr.Sel.Name) {
			return true
		}

		translation, valid := f.processMethod(selectorExpr.Sel.Name, callExpr)
		if valid {
			f.Translations = append(f.Translations, translation)
		}
		return true
	})

	f.Translations = cleanDuplicates(f.Translations)
}

func methodToIndex(method string) int {
	switch method {
	case Get:
		return 0
	case GetD:
		return 1
	default:
		return -1
	}
}
