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

type PoCompiler struct {
	File   *po.File // The source file containing translation entries.
	Config PoConfig // Configuration settings for compilation.
}

func NewPo(file *po.File, options ...PoOption) PoCompiler {
	return PoCompiler{
		File:   file,
		Config: DefaultPoConfig(options...),
	}
}

func (c PoCompiler) ToWriter(w io.Writer) error {
	var err error

	if c.Config.Verbose {
		c.Config.Logger.Println("Writing header...")
	}
	// Write the PO file header.
	err = c.writeHeader(w)
	if err != nil && !c.Config.IgnoreErrors {
		err = fmt.Errorf("error writing header format: %w", err)
		c.Config.Logger.Println("ERROR:", err)
	}

	if c.Config.Verbose {
		c.Config.Logger.Println("Cleaning duplicates...")
	}
	// Remove duplicate entries and write each entry to the writer.
	entries := c.File.Entries.CleanDuplicates().SortByFuzzy()
	if c.Config.Verbose {
		c.Config.Logger.Println("Writing entries...")
	}
	for i, e := range entries {
		if e.ID == "" {
			continue
		}
		err = c.writeEntry(w, e)
		if err != nil && !c.Config.IgnoreErrors {
			err = fmt.Errorf("error writing entry[%d]: %w", i, err)
			c.Config.Logger.Println("ERROR:", err)
			return err
		}
	}
	return nil
}

func (c PoCompiler) ToFile(f string) error {
	flags := os.O_WRONLY | os.O_TRUNC | os.O_CREATE
	if !c.Config.ForcePo {
		flags |= os.O_EXCL
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

func (c PoCompiler) ToString() string {
	var b strings.Builder

	// Write the compiled content to the string builder.
	c.ToWriter(&b)

	return b.String()
}

func (c PoCompiler) ToBytes() []byte {
	var b bytes.Buffer

	// Write the compiled content to the byte buffer.
	c.ToWriter(&b)

	return b.Bytes()
}
