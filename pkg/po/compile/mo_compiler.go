package compile

import (
	"bufio"
	"bytes"
	"io"
	"os"

	"github.com/Tom5521/gotext-tools/v2/pkg/po"
)

var _ po.Compiler = (*MoCompiler)(nil)

type MoCompiler struct {
	File   *po.File
	Config MoConfig
}

func NewMo(file *po.File, opts ...MoOption) MoCompiler {
	c := MoCompiler{
		File:   file,
		Config: DefaultMoConfig(opts...),
	}

	return c
}

func (mc *MoCompiler) SetFile(f *po.File) {
	mc.File = f
}

func (mc *MoCompiler) ToWriterWithOptions(w io.Writer, opts ...MoOption) error {
	mc.Config.ApplyOptions(opts...)
	defer mc.Config.RestoreLastCfg()
	return mc.ToWriter(w)
}

func (mc *MoCompiler) ToBytesWithOptions(options ...MoOption) []byte {
	mc.Config.ApplyOptions(options...)
	defer mc.Config.RestoreLastCfg()
	return mc.ToBytes()
}

func (mc *MoCompiler) ToFileWithOptions(f string, options ...MoOption) error {
	mc.Config.ApplyOptions(options...)
	defer mc.Config.RestoreLastCfg()
	return mc.ToFile(f)
}

func (mc MoCompiler) ToWriter(w io.Writer) error {
	buf := bufio.NewWriter(w)
	err := mc.writeTo(buf)

	if err != nil && !mc.Config.IgnoreErrors {
		return mc.error("error writing to buffer: %w", err)
	}

	mc.info("writing...")
	err = buf.Flush()
	if err != nil && !mc.Config.IgnoreErrors {
		return mc.error("error flushing buffer: %w", err)
	}

	return nil
}

func (mc MoCompiler) ToFile(f string) error {
	mc.info("oppening file...")
	// Open the file with the determined flags.
	flags := os.O_WRONLY | os.O_TRUNC | os.O_CREATE
	if !mc.Config.Force {
		flags |= os.O_EXCL
	}
	file, err := os.OpenFile(f, flags, os.ModePerm)
	if err != nil && !mc.Config.IgnoreErrors {
		err = mc.error("error opening file: %w", err)
		return err
	}
	defer file.Close()

	mc.info("cleaning file contents...")

	// Write compiled translations to the file.
	return mc.ToWriter(file)
}

func (mc MoCompiler) ToBytes() []byte {
	var b bytes.Buffer

	mc.ToWriter(&b)

	return b.Bytes()
}
