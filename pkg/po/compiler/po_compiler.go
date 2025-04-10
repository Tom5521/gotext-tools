package compiler

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/Tom5521/gotext-tools/pkg/po"
)

var _ po.Compiler = (*PoCompiler)(nil)

type PoCompiler struct {
	File   *po.File // The source file containing translation entries.
	Config PoConfig // Configuration settings for compilation.

	nplurals uint
	header   po.Header
}

func NewPo(file *po.File, options ...PoOption) PoCompiler {
	return PoCompiler{
		File:   file,
		Config: DefaultPoConfig(options...),
	}
}

func (c *PoCompiler) SetFile(f *po.File) {
	c.File = f
}

func (c *PoCompiler) ToWriterWithOptions(w io.Writer, opts ...PoOption) error {
	c.Config.ApplyOptions(opts...)
	defer c.Config.RestoreLastCfg()
	return c.ToWriter(w)
}

func (c *PoCompiler) ToStringWithOptions(opts ...PoOption) string {
	c.Config.ApplyOptions(opts...)
	defer c.Config.RestoreLastCfg()
	return c.ToString()
}

func (c *PoCompiler) ToFileWithOptions(f string, opts ...PoOption) error {
	c.Config.ApplyOptions(opts...)
	defer c.Config.RestoreLastCfg()
	return c.ToFile(f)
}

func (c *PoCompiler) ToBytesWithOptions(opts ...PoOption) []byte {
	c.Config.ApplyOptions(opts...)
	defer c.Config.RestoreLastCfg()
	return c.ToBytes()
}

func (c *PoCompiler) init() {
	c.header = c.File.Header()
	c.nplurals = c.header.Nplurals()
}

func (c PoCompiler) ToWriter(w io.Writer) error {
	c.init()

	buf := bufio.NewWriter(w)

	if c.Config.Verbose {
		c.Config.Logger.Println("Writing header...")
	}

	c.writeHeader(buf)

	if c.Config.Verbose {
		c.Config.Logger.Println("Cleaning duplicates...")
	}
	// Remove duplicate entries and write each entry to the writer.
	entries := c.File.Entries
	if c.Config.Verbose {
		c.Config.Logger.Println("Writing entries...")
	}

	entries = slices.DeleteFunc(entries, func(e po.Entry) bool {
		return e.Context == "" && e.ID == ""
	})

	for _, e := range entries {
		c.writeEntry(buf, e)
	}

	err := buf.Flush()
	if err != nil && !c.Config.IgnoreErrors {
		return err
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
