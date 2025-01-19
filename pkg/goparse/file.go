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

// File represents a Go source file that is being processed by the parser.
// This struct contains information about the file's content, path, package,
// and whether it imports the desired package for translation processing.
type File struct {
	// Config is a pointer to the configuration used for parsing.
	// Using a pointer avoids excessive memory usage when working with many files
	// and allows changes to the parser configuration to propagate to each file.
	config *poconfig.Config

	file      *ast.File // The parsed abstract syntax tree (AST) of the file.
	content   []byte    // The raw content of the file as a byte slice.
	path      string    // The path to the file.
	pkgName   string    // The name of the package declared in the file.
	hasGotext bool      // Indicates if the file imports the desired "gotext" package.
}

// WantedImport specifies the import path of the "gotext" library
// that the parser is looking for in the Go files.
const WantedImport = `"github.com/leonelquinteros/gotext"`

// NewFileFromReader creates a new File instance by reading content from an io.Reader.
// The content is read into memory and processed according to the provided configuration.
func NewFileFromReader(r io.Reader, name string, cfg *poconfig.Config) (*File, error) {
	err := validateConfig(*cfg)
	if err != nil {
		return nil, err
	}
	return unsafeNewFileFromReader(r, name, cfg)
}

// unsafeNewFileFromReader is an internal method that creates a File instance
// from an io.Reader without validating the configuration.
func unsafeNewFileFromReader(r io.Reader, name string, cfg *poconfig.Config) (*File, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return unsafeNewFile(content, name, cfg)
}

// NewFileFromPath creates a new File instance by reading content from a file on disk.
func NewFileFromPath(path string, cfg *poconfig.Config) (*File, error) {
	err := validateConfig(*cfg)
	if err != nil {
		return nil, err
	}

	return unsafeNewFileFromPath(path, cfg)
}

// unsafeNewFileFromPath is an internal method that creates a File instance
// from a file on disk without validating the configuration.
func unsafeNewFileFromPath(path string, cfg *poconfig.Config) (*File, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return unsafeNewFileFromReader(file, path, cfg)
}

// NewFileFromBytes creates a new File instance from raw byte data.
func NewFileFromBytes(b []byte, name string, cfg *poconfig.Config) (*File, error) {
	err := validateConfig(*cfg)
	if err != nil {
		return nil, err
	}

	return unsafeNewFile(b, name, cfg)
}

// unsafeNewFile is an internal method that creates a File instance from raw byte data
// and the provided configuration without validating the configuration.
func unsafeNewFile(content []byte, name string, cfg *poconfig.Config) (*File, error) {
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

// determinePackageInfo analyzes the file's AST to extract package-related information.
// It determines the package name and checks if the desired "gotext" package is imported.
func (f *File) determinePackageInfo() {
	f.pkgName = "gotext" // Default package name.
	for _, imprt := range f.file.Imports {
		f.hasGotext = imprt.Path.Value == WantedImport
		if f.hasGotext {
			// If the import is aliased, use the alias as the package name.
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
