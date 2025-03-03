package compiler

import (
	"bytes"
	bin "encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unsafe"

	"github.com/Tom5521/xgotext/pkg/po"
)

var _ Compiler = (*MoCompiler)(nil)

var (
	isBigEndian = (*(*[2]uint8)(unsafe.Pointer(&[]uint16{1}[0])))[0] == 0

	magicNumber = func() uint32 {
		if isBigEndian {
			return bigEndianMagicNumber
		}
		return littleEndianMagicNumber
	}()
	order        = bin.NativeEndian
	eotSeparator = string([]byte{0x4})
	nulSeparator = string([]byte{0x0})
)

const (
	bigEndianMagicNumber    uint32 = 0xde120495
	littleEndianMagicNumber uint32 = 0x950412de
)

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

func (mc MoCompiler) createBinary() ([]byte, error) {
	entries := mc.File.Entries.CleanDuplicates().Solve()

	var offsets []int
	var ids, strs strings.Builder
	for _, e := range entries {
		var msgid string
		var msgstr string
		if e.Context != "" {
			msgid = e.Context + eotSeparator
		}
		if e.Plural != "" {
			var msgstrs []string
			plurals := e.Plurals.Sort()
			for _, plural := range plurals {
				msgstrs = append(msgstrs, plural.Str)
			}
			msgid += e.ID + nulSeparator + e.Plural
			msgstr = strings.Join(msgstrs, nulSeparator)
		} else {
			msgid += e.ID
			msgstr = e.Str
		}

		offsets = append(offsets,
			ids.Len(),
			len(msgid),
			strs.Len(),
			len(msgstr),
		)
		ids.WriteString(msgid + nulSeparator)
		strs.WriteString(msgstr + nulSeparator)
	}

	keystart := 7*4 + 16*len(entries)
	valuestart := keystart + ids.Len()

	var koffsets, voffsets []int

	for i := 0; i < len(offsets); i += 4 {
		if i+3 >= len(offsets) {
			return nil, errors.New("not enough values to unpack")
		}
		o1 := offsets[i]
		l1 := offsets[i+1]
		o2 := offsets[i+2]
		l2 := offsets[i+3]

		koffsets = append(koffsets, l1, o1+keystart)
		voffsets = append(voffsets, l2, o2+valuestart)
	}
	offsets = append(koffsets, voffsets...)

	// bytes alias
	bts := make([]byte, len(offsets)*4)
	for i, v := range offsets {
		order.AppendUint32(bts[i*4:], uint32(v))
	}

	var (
		err  error
		buf  bytes.Buffer
		data = []any{
			magicNumber,
			0,
			len(entries),
			7 * 4,
			7*4 + len(entries)*8,
			0, keystart,
			bts,
			// ids.String(),
			// strs.String(),
			[]byte(ids.String()),
			[]byte(strs.String()),
		}
	)

	// Fix not fixed-size integers.
	for i, v := range data {
		if _, ok := v.(int); ok {
			data[i] = int32(v.(int))
		}
	}

	// Write data.
	for _, value := range data {
		err = bin.Write(&buf, order, value)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func (mc MoCompiler) ToWriter(w io.Writer) error {
	b, err := mc.createBinary()
	if err != nil && !mc.Config.IgnoreErrors {
		return err
	}

	_, err = w.Write(b)
	if err != nil && !mc.Config.IgnoreErrors {
		return err
	}

	return nil
}

func (c MoCompiler) ToFile(f string) error {
	if c.Config.Verbose {
		c.Config.Logger.Println("Opening file...")
	}
	// Open the file with the determined flags.
	file, err := os.OpenFile(f, os.O_RDWR, os.ModePerm)
	if err != nil && !c.Config.IgnoreErrors {
		err = fmt.Errorf("error opening file: %w", err)
		c.Config.Logger.Println("ERROR:", err)
		return err
	}
	defer file.Close()

	if c.Config.Verbose {
		c.Config.Logger.Println("Cleaning file contents...")
	}
	file, err = os.Create(file.Name())
	if err != nil {
		err = fmt.Errorf("error truncating file: %w", err)
		c.Config.Logger.Println("ERROR:", err)
		return err
	}

	// Write compiled translations to the file.
	return c.ToWriter(file)
}

func (mc MoCompiler) ToString() string {
	var b strings.Builder

	bin, _ := mc.createBinary()
	b.Write(bin)

	return b.String()
}

func (mc MoCompiler) ToBytes() []byte {
	b, _ := mc.createBinary()

	return b
}
