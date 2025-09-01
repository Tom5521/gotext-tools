package compile

import (
	"bufio"
	"bytes"
	"io"
	"os"

	"github.com/Tom5521/gotext-tools/v2/pkg/po"
)

var _ po.Compiler = (*MoCompiler)(nil)

// MoCompiler implements the po.Compiler interface for compiling PO files to MO format.
// It holds a PO file reference and configuration settings for the compilation process.
type MoCompiler struct {
	// File is the PO file to be compiled
	File *po.File
	// Config contains the compilation configuration
	Config MoConfig
}

// NewMo creates a new MoCompiler instance with the given PO file and optional configuration.
// It applies any provided MoOption functions to configure the compiler.
func NewMo(file *po.File, opts ...MoOption) MoCompiler {
	c := MoCompiler{
		File:   file,
		Config: DefaultMoConfig(opts...),
	}

	return c
}

// SetFile updates the PO file reference in the compiler.
func (mc *MoCompiler) SetFile(f *po.File) {
	mc.File = f
}

// ToWriterWithOptions writes the compiled MO output to an io.Writer with temporary options.
// The options are only applied for this operation and then reverted.
// Returns an error if compilation or writing fails.
func (mc *MoCompiler) ToWriterWithOptions(w io.Writer, opts ...MoOption) error {
	mc.Config.ApplyOptions(opts...)
	defer mc.Config.RestoreLastCfg()
	return mc.ToWriter(w)
}

// ToBytesWithOptions returns the compiled MO data as a byte slice with temporary options.
// The options are only applied for this operation and then reverted.
func (mc *MoCompiler) ToBytesWithOptions(options ...MoOption) []byte {
	mc.Config.ApplyOptions(options...)
	defer mc.Config.RestoreLastCfg()
	return mc.ToBytes()
}

// ToFileWithOptions writes the compiled MO output to a file with temporary options.
// The options are only applied for this operation and then reverted.
// Returns an error if file operations or compilation fails.
func (mc *MoCompiler) ToFileWithOptions(f string, options ...MoOption) error {
	mc.Config.ApplyOptions(options...)
	defer mc.Config.RestoreLastCfg()
	return mc.ToFile(f)
}

// ToWriter writes the compiled MO output to an io.Writer.
// It uses buffered writing for better performance with large files.
// Returns an error if compilation or writing fails.
func (mc MoCompiler) ToWriter(w io.Writer) error {
	buf := bufio.NewWriter(w)
	err := mc.writeTo(buf)
	if err != nil {
		return mc.error("error writing to buffer: %w", err)
	}

	mc.info("writing...")
	err = buf.Flush()
	if err != nil {
		return mc.error("error flushing buffer: %w", err)
	}

	return nil
}

// ToFile writes the compiled MO output to the specified file path.
// By default, it fails if the file already exists (unless Force is enabled).
// Returns an error if file operations or compilation fails.
func (mc MoCompiler) ToFile(f string) error {
	// Open the file with the determined flags.
	flags := os.O_WRONLY | os.O_TRUNC | os.O_CREATE
	if !mc.Config.Force {
		flags |= os.O_EXCL
	}
	mc.info("opening file...")
	file, err := os.OpenFile(f, flags, 0o600)
	if err != nil && !mc.Config.IgnoreErrors {
		err = mc.error("error opening file: %w", err)
		return err
	}
	defer file.Close()

	// Write compiled translations to the file.
	return mc.ToWriter(file)
}

// ToBytes returns the compiled MO data as a byte slice.
// This is useful for in-memory operations or when you need the raw binary data.
func (mc MoCompiler) ToBytes() []byte {
	var b bytes.Buffer

	mc.ToWriter(&b)

	return b.Bytes()
}
