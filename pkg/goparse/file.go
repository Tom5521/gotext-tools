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
	"github.com/Tom5521/xgotext/pkg/po/config"
	"github.com/Tom5521/xgotext/pkg/po/entry"
)

// Constants.
const (
	DefaultPackageName = "gotext"
	// WantedImport specifies the import path of the "gotext" library
	// that the parser is looking for in the Go files.
	WantedImport = `"github.com/leonelquinteros/gotext"`
)

// File represents a Go source file that is being processed by the parser.
// This struct contains information about the file's content, path, package,
// and whether it imports the desired package for translation processing.
type File struct {
	// Config is a pointer to the configuration used for parsing.
	// Using a pointer avoids excessive memory usage when working with many files
	// and allows changes to the parser configuration to propagate to each file.
	config     *config.Config
	seenTokens map[ast.Node]bool
	file       *ast.File // The parsed abstract syntax tree (AST) of the file.
	content    []byte    // The raw content of the file as a byte slice.
	path       string    // The path to the file.
	pkgName    string    // The name of the package declared in the file.
	hasGotext  bool      // Indicates if the file imports the desired "gotext" package.
}

// NewFileFromReader creates a new File instance by reading content from an io.Reader.
// The content is read into memory and processed according to the provided configuration.
func NewFileFromReader(r io.Reader, name string, cfg *config.Config) (*File, error) {
	err := validateConfig(*cfg)
	if err != nil {
		return nil, err
	}
	return unsafeNewFileFromReader(r, name, cfg)
}

// unsafeNewFileFromReader is an internal method that creates a File instance
// from an io.Reader without validating the configuration.
func unsafeNewFileFromReader(r io.Reader, name string, cfg *config.Config) (*File, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read content: %w", err)
	}

	return unsafeNewFile(content, name, cfg)
}

// NewFileFromPath creates a new File instance by reading content from a file on disk.
func NewFileFromPath(path string, cfg *config.Config) (*File, error) {
	err := validateConfig(*cfg)
	if err != nil {
		return nil, err
	}

	return unsafeNewFileFromPath(path, cfg)
}

// unsafeNewFileFromPath is an internal method that creates a File instance
// from a file on disk without validating the configuration.
func unsafeNewFileFromPath(path string, cfg *config.Config) (*File, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	return unsafeNewFileFromReader(file, path, cfg)
}

// NewFileFromBytes creates a new File instance from raw byte data.
func NewFileFromBytes(b []byte, name string, cfg *config.Config) (*File, error) {
	err := validateConfig(*cfg)
	if err != nil {
		return nil, err
	}

	return unsafeNewFile(b, name, cfg)
}

// unsafeNewFile is an internal method that creates a File instance from raw byte data
// and the provided configuration without validating the configuration.
func unsafeNewFile(content []byte, name string, cfg *config.Config) (*File, error) {
	file := &File{
		content:    content,
		path:       name,
		config:     cfg,
		seenTokens: make(map[ast.Node]bool),
		pkgName:    DefaultPackageName,
	}

	if err := file.parse(); err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	file.determinePackageInfo()
	return file, nil
}

// parse parses the file content into an AST
func (f *File) parse() error {
	parsedFile, err := parser.ParseFile(token.NewFileSet(), f.path, f.content, 0)
	if err != nil {
		return err
	}
	f.file = parsedFile
	return nil
}

// determinePackageInfo analyzes the file's AST to extract package-related information.
// It determines the package name and checks if the desired "gotext" package is imported.
func (f *File) determinePackageInfo() {
	for _, imp := range f.file.Imports {
		if imp.Path.Value == WantedImport {
			f.hasGotext = true
			if imp.Name != nil {
				f.pkgName = imp.Name.String()
			}
			break
		}
	}
}

// isGotextCall checks if an AST node represents a gotext function call
func (f *File) isGotextCall(n ast.Node) bool {
	callExpr, ok := n.(*ast.CallExpr)
	if !ok {
		return false
	}

	selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident, ok := selectorExpr.X.(*ast.Ident)
	if !ok || ident.Name != f.pkgName {
		return false
	}

	_, ok = translationMethods[selectorExpr.Sel.Name]
	return ok
}

func (f *File) basicLitToTranslation(n *ast.BasicLit) (entry.Translation, error) {
	str, err := strconv.Unquote(n.Value)
	if err != nil {
		return entry.Translation{}, err
	}

	return entry.Translation{
		ID: str,
		Locations: []entry.Location{{
			Line: util.FindLine(f.content, n.Pos()),
			File: f.path,
		}},
	}, nil
}

func (f *File) extractCommonString(n ast.Node) ([]entry.Translation, []error) {
	var translations []entry.Translation
	var errors []error

	processExpressions := func(exprs []ast.Expr) {
		for _, expr := range exprs {
			if lit, ok := expr.(*ast.BasicLit); ok && lit.Kind == token.STRING {
				if f.seenTokens[lit] {
					continue
				}

				if lit.Value == `""` {
					continue
				}

				translation, err := f.basicLitToTranslation(lit)
				if err != nil {
					errors = append(errors, err)
					continue
				}

				translations = append(translations, translation)
				f.seenTokens[lit] = true
			}
		}
	}

	switch t := n.(type) {
	case *ast.CallExpr:
		processExpressions(t.Args)
	case *ast.AssignStmt:
		processExpressions(t.Rhs)
	case *ast.ValueSpec:
		processExpressions(t.Values)
	}

	return translations, errors
}

// Translations returns all translations found in the file
func (f *File) Translations() ([]entry.Translation, []error) {
	var translations []entry.Translation
	var errors []error

	for n := range InspectNode(f.file) {
		if n == nil {
			continue
		}

		if !f.isGotextCall(n) {
			if f.config.ExtractAll {
				ts, errs := f.extractCommonString(n)
				translations = append(translations, ts...)
				errors = append(errors, errs...)
			}
			continue
		}

		callExpr := n.(*ast.CallExpr)
		selectorExpr := callExpr.Fun.(*ast.SelectorExpr)

		translation, valid, err := f.processMethod(selectorExpr.Sel.Name, callExpr)
		if err != nil {
			errors = append(errors, fmt.Errorf("[%s:%d] %w",
				f.path,
				util.FindLine(f.content, selectorExpr.Pos()),
				err))
		}
		if valid && err == nil {
			translations = append(translations, translation)
		}
	}

	return util.CleanDuplicates(translations), errors
}
