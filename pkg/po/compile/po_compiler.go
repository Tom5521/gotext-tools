package compile

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Tom5521/gotext-tools/v2/pkg/po"
)

var _ po.Compiler = (*PoCompiler)(nil)

type PoCompiler struct {
	File   *po.File // The source file containing translation entries.
	Config PoConfig // Configuration settings for compilation.

	nplurals uint
	header   po.Header
}

func (p PoCompiler) error(format string, a ...any) error {
	if p.Config.IgnoreErrors {
		return nil
	}

	format = "compile: " + format
	err := fmt.Errorf(format, a...)
	if p.Config.Logger != nil {
		p.Config.Logger.Println("ERROR:", err)
	}

	return err
}

func (p PoCompiler) info(format string, a ...any) {
	if p.Config.Logger != nil && p.Config.Verbose {
		p.Config.Logger.Println("INFO:", fmt.Sprintf(format, a...))
	}
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

	var reader bytes.Buffer
	buf := bufio.NewWriter(w)
	var writer io.Writer = buf
	if c.Config.Highlight != nil {
		writer = io.MultiWriter(buf, &reader)
	}

	c.info("writing header...")
	c.writeHeader(writer)
	c.info("cleaning duplicates...")
	// Remove duplicate entries and write each entry to the writer.
	entries := c.File.CutHeader().CleanDuplicates()
	c.info("writing entries...")

	for _, e := range entries {
		c.writeEntry(writer, e)
	}

	if c.Config.Highlight != nil {
		c.info("highlighting info...")
		h, err := HighlightFromBytes(
			c.Config.Highlight,
			c.File.Name,
			reader.Bytes(),
		)
		if err != nil {
			c.error("error highlighting output: %w", err)
		}
		buf.Reset(w)
		buf.Write(h)
	}

	err := buf.Flush()
	if err != nil {
		return c.error("error flushing buffer: %w", err)
	}

	return nil
}

func (c PoCompiler) ToFile(f string) error {
	flags := os.O_WRONLY | os.O_TRUNC | os.O_CREATE
	if !c.Config.ForcePo {
		flags |= os.O_EXCL
	}

	c.info("opening file...")

	// Open the file with the determined flags.
	file, err := os.OpenFile(f, flags, os.ModePerm)
	if err != nil {
		return c.error("error opening file: %w", err)
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
