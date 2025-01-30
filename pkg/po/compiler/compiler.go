package compiler

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Tom5521/xgotext/pkg/po/config"
	"github.com/Tom5521/xgotext/pkg/po/types"
)

type Compiler struct {
	File   *types.File
	Config config.Config
}

// CompileToWriter writes the compiled translations to an `io.Writer` in the PO file format.
func (c Compiler) CompileToWriter(w io.Writer) error {
	var err error

	fmt.Fprintln(w, c.formatHeader())

	// Clean duplicates in translations and write each to the writer.
	translations := types.CleanDuplicates(c.File.Entries)
	for _, t := range translations {
		if t.ID == "" {
			continue
		}
		_, err = fmt.Fprintln(w, c.formatEntry(t))
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
