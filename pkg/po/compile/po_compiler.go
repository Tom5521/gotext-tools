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
	// File is the source PO file containing translation entries
	File *po.File
	// Config contains compilation settings and options
	Config PoConfig

	nplurals uint      // Number of plural forms from the header
	header   po.Header // Parsed header information

	writer io.Writer
}

// error creates and logs an error message if error reporting is enabled.
// Returns nil if IgnoreErrors is true.
func (c PoCompiler) error(format string, a ...any) error {
	if c.Config.IgnoreErrors {
		return nil
	}

	format = "compile: " + format
	err := fmt.Errorf(format, a...)
	if c.Config.Logger != nil {
		c.Config.Logger.Println("ERROR:", err)
	}

	return err
}

// info logs an informational message if verbose logging is enabled.
func (c PoCompiler) info(format string, a ...any) {
	if c.Config.Logger != nil && c.Config.Verbose {
		c.Config.Logger.Println("INFO:", fmt.Sprintf(format, a...))
	}
}

// NewPo creates a new PoCompiler instance with the given PO file and options.
// The options are applied to configure the compiler behavior.
func NewPo(file *po.File, options ...PoOption) PoCompiler {
	return PoCompiler{
		File:   file,
		Config: DefaultPoConfig(options...),
	}
}

// SetFile updates the PO file reference in the compiler.
func (c *PoCompiler) SetFile(f *po.File) {
	c.File = f
}

// ToWriterWithOptions writes compiled output to an io.Writer with temporary options.
// The options are only applied for this operation and then reverted.
func (c *PoCompiler) ToWriterWithOptions(w io.Writer, opts ...PoOption) error {
	c.Config.ApplyOptions(opts...)
	defer c.Config.RestoreLastCfg()
	return c.ToWriter(w)
}

// ToStringWithOptions returns the compiled output as a string with temporary options.
// The options are only applied for this operation and then reverted.
func (c *PoCompiler) ToStringWithOptions(opts ...PoOption) string {
	c.Config.ApplyOptions(opts...)
	defer c.Config.RestoreLastCfg()
	return c.ToString()
}

// ToFileWithOptions writes compiled output to a file with temporary options.
// The options are only applied for this operation and then reverted.
func (c *PoCompiler) ToFileWithOptions(f string, opts ...PoOption) error {
	c.Config.ApplyOptions(opts...)
	defer c.Config.RestoreLastCfg()
	return c.ToFile(f)
}

// ToBytesWithOptions returns the compiled output as bytes with temporary options.
// The options are only applied for this operation and then reverted.
func (c *PoCompiler) ToBytesWithOptions(opts ...PoOption) []byte {
	c.Config.ApplyOptions(opts...)
	defer c.Config.RestoreLastCfg()
	return c.ToBytes()
}

// init initializes the compiler by parsing header information.
func (c *PoCompiler) init() {
	if c.Config.ManageHeader {
		c.header = c.File.Header()
	}
	c.nplurals = c.header.Nplurals()
}

// ToWriter writes compiled PO content to an io.Writer.
// Handles header writing, duplicate cleaning, and optional syntax highlighting.
func (c PoCompiler) ToWriter(w io.Writer) error {
	c.init()

	var reader bytes.Buffer
	buf := bufio.NewWriter(w)
	var writer io.Writer = buf
	if c.Config.Highlight != nil {
		writer = io.MultiWriter(buf, &reader)
	}

	if c.Config.ManageHeader {
		c.info("writing header...")
		if !c.Config.OmitHeader {
			if c.Config.HeaderConfig != nil {
				c.header = c.Config.HeaderConfig.ToHeader()
			}
			c.writeHeader(writer)
		}
	}
	entries := c.File.Entries

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
		// Second revision: Why the hell does it send an error
		// but the algorithm continues regardless, and why
		// does it still write h to the buffer?
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

// ToFile writes compiled output to the specified file path.
// By default fails if file exists (unless ForcePo is enabled).
func (c PoCompiler) ToFile(f string) error {
	flags := os.O_WRONLY | os.O_TRUNC | os.O_CREATE
	if !c.Config.ForcePo {
		flags |= os.O_EXCL
	}

	c.info("opening file...")
	file, err := os.OpenFile(f, flags, os.ModePerm)
	if err != nil {
		return c.error("error opening file: %w", err)
	}
	defer file.Close()

	return c.ToWriter(file)
}

// ToString returns the compiled output as a string.
func (c PoCompiler) ToString() string {
	var b strings.Builder
	c.ToWriter(&b)
	return b.String()
}

// ToBytes returns the compiled output as a byte slice.
func (c PoCompiler) ToBytes() []byte {
	var b bytes.Buffer
	c.ToWriter(&b)
	return b.Bytes()
}
