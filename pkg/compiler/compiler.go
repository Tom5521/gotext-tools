package compiler

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/poconfig"
	"github.com/Tom5521/xgotext/pkg/poentry"
)

//go:embed header.pot
var baseHeader string

// Compiler is responsible for compiling a list of translations into various formats
// (e.g., string, file, or bytes) based on the given configuration.
type Compiler struct {
	Translations []poentry.Translation // List of translations to compile.
	Config       poconfig.Config       // Configuration for the compilation process.
}

// CompileToWriter writes the compiled translations to an `io.Writer` in the PO file format.
func (c Compiler) CompileToWriter(w io.Writer) error {
	// Write the base header, including package version, language, and plural forms.
	_, err := fmt.Fprintf(
		w,
		baseHeader,
		c.Config.PackageVersion,
		c.Config.Language,
		c.Config.Nplurals,
	)
	if err != nil {
		return err
	}

	// Clean duplicates in translations and write each to the writer.
	translations := util.CleanDuplicates(c.Translations)
	for _, t := range translations {
		_, err = fmt.Fprintln(w, t.Format(c.Config))
		if err != nil {
			return err
		}
	}
	return nil
}

// CompileToFile writes the compiled translations to a specified file. If `ForcePo` is set in the configuration,
// the file is created or truncated before writing.
func (c Compiler) CompileToFile(f string) error {
	flags := os.O_RDWR
	if c.Config.ForcePo {
		flags |= os.O_CREATE
	}
	file, err := os.OpenFile(f, flags, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	// Truncate the file if ForcePo is enabled.
	if c.Config.ForcePo {
		err = file.Truncate(0)
		if err != nil {
			return err
		}
	}

	return c.CompileToWriter(file)
}

// CompileToString compiles the translations and returns the result as a string.
func (c Compiler) CompileToString() string {
	var b strings.Builder

	c.CompileToWriter(&b)

	return b.String()
}

// CompileToBytes compiles the translations and returns the result as a byte slice.
func (c Compiler) CompileToBytes() []byte {
	var b bytes.Buffer

	c.CompileToWriter(&b)

	return b.Bytes()
}

// ErrNotImplementedYet is used as an error for functions that are not yet implemented.
var ErrNotImplementedYet = errors.New("not implemented yet (sorry)")

// CompileToDir compiles the translations to a directory. This function is not implemented yet.
func (c Compiler) CompileToDir(d string) error { return ErrNotImplementedYet }

// CompileToDirFS compiles the translations to a directory in a filesystem. This function is not implemented yet.
func (c Compiler) CompileToDirFS(dfs fs.FS) error { return ErrNotImplementedYet }
