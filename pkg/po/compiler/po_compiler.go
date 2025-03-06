package compiler

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Tom5521/xgotext/pkg/po"
)

var _ Compiler = (*PoCompiler)(nil)

// PoCompiler is responsible for compiling translations from a `types.File`
// into different output formats, such as strings, byte slices, or files.
type PoCompiler struct {
	File   *po.File // The source file containing translation entries.
	Config PoConfig // Configuration settings for compilation.
}

// NewPo creates a new Compiler instance with the given translation file and options.
// The provided options override the default configuration.
func NewPo(file *po.File, options ...PoOption) PoCompiler {
	return PoCompiler{
		File:   file,
		Config: DefaultPoConfig(options...),
	}
}

// ToWriter writes the compiled translations to an `io.Writer` in the PO file format.
// The provided options override the instance's configuration.
func (c PoCompiler) ToWriter(w io.Writer) error {
	// Apply the provided options, which take precedence over the instance's configuration.
	var err error

	if c.Config.Verbose {
		c.Config.Logger.Println("Writing header...")
	}
	// Write the PO file header.
	_, err = fmt.Fprintln(w, c.formatHeader())
	if err != nil && !c.Config.IgnoreErrors {
		err = fmt.Errorf("error writing header format: %w", err)
		c.Config.Logger.Println("ERROR:", err)
		return err
	}
	if c.Config.Verbose {
		c.Config.Logger.Println("Cleaning duplicates...")
	}
	// Remove duplicate entries and write each entry to the writer.
	entries := c.File.Entries.CleanDuplicates()
	if c.Config.Verbose {
		c.Config.Logger.Println("Writing entries...")
	}
	for i, e := range entries {
		if e.ID == "" {
			continue
		}

		_, err = fmt.Fprintln(w, c.formatEntry(e))
		if err != nil && !c.Config.IgnoreErrors {
			err = fmt.Errorf("error writing entry[%d]: %w", i, err)
			c.Config.Logger.Println("ERROR:", err)
			return err
		}
	}
	return nil
}

// ToFile writes the compiled translations to a specified file.
// If `ForcePo` is enabled, the file is created or truncated before writing.
// The provided options override the instance's configuration.
func (c PoCompiler) ToFile(f string) error {
	flags := os.O_WRONLY | os.O_TRUNC
	if c.Config.ForcePo {
		flags |= os.O_CREATE
	}

	if c.Config.Verbose {
		c.Config.Logger.Println("Opening file...")
	}

	// Open the file with the determined flags.
	file, err := os.OpenFile(f, flags, os.ModePerm)
	if err != nil && !c.Config.IgnoreErrors {
		err = fmt.Errorf("error opening file: %w", err)
		c.Config.Logger.Println("ERROR:", err)
		return err
	}
	defer file.Close()

	// Write compiled translations to the file.
	return c.ToWriter(file)
}

// ToString compiles the translations and returns the result as a string.
// The provided options override the instance's configuration.
func (c PoCompiler) ToString() string {
	var b strings.Builder

	// Write the compiled content to the string builder.
	c.ToWriter(&b)

	return b.String()
}

// ToBytes compiles the translations and returns the result as a byte slice.
// The provided options override the instance's configuration.
func (c PoCompiler) ToBytes() []byte {
	var b bytes.Buffer

	// Write the compiled content to the byte buffer.
	c.ToWriter(&b)

	return b.Bytes()
}
