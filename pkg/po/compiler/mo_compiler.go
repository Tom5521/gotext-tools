package compiler

import (
	"bytes"
	bin "encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
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
	eotSeparator = "\x04"
	nulSeparator = "\x00"
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

func cleanEntries(in po.Entries) (out po.Entries) {
	in = in.Solve()

	for _, v := range in {
		if slices.Contains(v.Flags, "fuzzy") {
			continue
		}
		out = append(out, v)
	}

	return
}

// Code translated from: https://github.com/izimobil/polib/blob/master/polib.py#L553
func (mc MoCompiler) writeTo(writer io.Writer) error {
	entries := cleanEntries(mc.File.Entries)

	var offsets []int
	var ids, strs string
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
			len(ids),
			len(msgid),
			len(strs),
			len(msgstr),
		)
		ids += msgid + nulSeparator
		strs += msgstr + nulSeparator
	}

	keystart := 7*4 + 16*len(entries)
	valuestart := keystart + len(ids)

	var koffsets, voffsets []int

	for i := 0; i < len(offsets); i += 4 {
		if i+3 >= len(offsets) {
			return errors.New("not enough values to unpack")
		}
		o1 := offsets[i]
		l1 := offsets[i+1]
		o2 := offsets[i+2]
		l2 := offsets[i+3]

		koffsets = append(koffsets, l1, o1+keystart)
		voffsets = append(voffsets, l2, o2+valuestart)
	}
	offsets = append(koffsets, voffsets...)

	data := []any{
		magicNumber,
		0,
		len(entries),
		7 * 4,
		7*4 + len(entries)*8,
		0, keystart,
		func() (s []uint32) {
			for _, v := range offsets {
				s = append(s, uint32(v))
			}
			return
		}(),
		[]byte(ids),
		[]byte(strs),
	}

	// Fix not fixed-size integers.
	for i, v := range data {
		if _, ok := v.(int); ok {
			data[i] = uint32(v.(int))
		}
	}

	// Write data.
	for _, value := range data {
		err := bin.Write(writer, order, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (mc MoCompiler) ToWriter(w io.Writer) error {
	err := mc.writeTo(w)
	if err != nil && !mc.Config.IgnoreErrors {
		return err
	}

	return nil
}

func (mc MoCompiler) ToFile(f string) error {
	if mc.Config.Verbose {
		mc.Config.Logger.Println("Opening file...")
	}
	// Open the file with the determined flags.
	flags := os.O_WRONLY | os.O_TRUNC
	if mc.Config.Force {
		flags |= os.O_CREATE
	}
	file, err := os.OpenFile(f, flags, os.ModePerm)
	if err != nil && !mc.Config.IgnoreErrors {
		err = fmt.Errorf("error opening file: %w", err)
		mc.Config.Logger.Println("ERROR:", err)
		return err
	}
	defer file.Close()

	if mc.Config.Verbose {
		mc.Config.Logger.Println("Cleaning file contents...")
	}

	// Write compiled translations to the file.
	return mc.ToWriter(file)
}

func (mc MoCompiler) ToBytes() []byte {
	var b bytes.Buffer

	mc.writeTo(&b)

	return b.Bytes()
}
