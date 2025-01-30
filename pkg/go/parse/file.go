package parse

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"

	"github.com/Tom5521/xgotext/pkg/po/config"
	"github.com/Tom5521/xgotext/pkg/po/types"
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

// parse parses the file content into an AST.
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

// Translations returns all translations found in the file.
func (f *File) Translations() ([]types.Entry, []error) {
	if f.config.Logger != nil && f.config.Verbose {
		f.config.Logger.Printf("Parsing %s...", f.path)
	}

	var translations []types.Entry
	var errors []error

	if !f.hasGotext {
		return translations, errors
	}

	for n := range InspectNode(f.file) {
		t, e := f.processNode(n)
		translations = append(translations, t...)
		errors = append(errors, e...)
	}

	return types.CleanDuplicates(translations), errors
}
