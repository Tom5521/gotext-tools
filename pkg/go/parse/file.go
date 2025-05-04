package parse

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"

	"github.com/Tom5521/gotext-tools/v2/pkg/po"
)

// Constants.
const (
	DefaultPackageName = "gotext"
	// WantedImport specifies the import path of the "gotext" library
	// that the parser is looking for in the Go files.
	WantedImport = `"github.com/leonelquinteros/gotext"`
)

// File represents the parser of an individual go file,
// I do it this way to keep track of the location and name of each file.
//
// It does not generate Header, it only extracts the entries according to the configuration.
type File struct {
	config    *Config
	seenNodes map[ast.Node]bool
	file      *ast.File // The parsed abstract syntax tree (AST) of the file.
	reader    *bytes.Reader
	name      string // The path to the file.
	pkgName   string // The name of the package declared in the file.
	hasGotext bool   // Indicates if the file imports the desired "gotext" package.

	errors []error
}

func (f *File) error(format string, a ...any) error {
	var err error
	format = "go/parse: " + format
	if len(a) == 0 {
		err = errors.New(format)
	} else {
		err = fmt.Errorf(format, a...)
	}

	if f.config.Logger != nil {
		f.config.Logger.Println("ERROR:", err)
	}

	f.errors = append(f.errors, err)

	return err
}

func (f *File) Reset(d io.Reader, name string, config *Config) error {
	f.seenNodes = nil
	f.errors = nil

	if r, ok := d.(*bytes.Reader); ok {
		f.reader = r
	} else {
		b, err := io.ReadAll(d)
		if err != nil {
			return err
		}
		f.reader = bytes.NewReader(b)
	}
	f.name = name

	if config != nil {
		f.config = config
	} else {
		f.config = &[]Config{DefaultConfig()}[0]
	}

	if err := f.parse(); err != nil {
		return err
	}

	f.determinePackageInfo()

	return nil
}

func NewFileFromBytes(b []byte, name string, config *Config) (*File, error) {
	file := &File{
		reader:  bytes.NewReader(b),
		name:    name,
		pkgName: DefaultPackageName,
		config:  config,
	}

	if err := file.parse(); err != nil {
		return nil, err
	}

	file.determinePackageInfo()
	return file, nil
}

// NewFileFromPath creates a new File instance by reading content from a file on disk.
func NewFileFromPath(path string, config *Config) (*File, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	return NewFile(file, path, config)
}

// NewFile creates a new File instance from raw byte data.
func NewFile(b io.Reader, name string, config *Config) (*File, error) {
	bytedata, err := io.ReadAll(b)
	if err != nil {
		return nil, err
	}

	return NewFileFromBytes(bytedata, name, config)
}

// parse parses the file content into an AST.
func (f *File) parse() error {
	var err error
	f.file, err = parser.ParseFile(token.NewFileSet(), f.name, f.reader, 0)
	if err != nil {
		return fmt.Errorf("failed to parse the file: %w", err)
	}
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

func (f *File) Errors() []error {
	return f.errors
}

func (f *File) Error() error {
	if len(f.errors) == 0 {
		return nil
	}

	return f.errors[0]
}

// Entries returns all translations found in the file.
func (f *File) Entries() po.Entries {
	// Reset fields.
	f.seenNodes = make(map[ast.Node]bool)
	f.errors = nil

	var entries po.Entries

	if !f.hasGotext && !f.config.ExtractAll {
		return entries
	}

	ast.Inspect(f.file, func(n ast.Node) bool {
		t, e := f.processNode(n)
		entries = append(entries, t...)
		f.errors = append(f.errors, e...)
		return true
	})

	return entries
}
