package compile

import (
	bin "encoding/binary"
	"fmt"
	"io"
	"reflect"

	"github.com/Tom5521/gotext-tools/v2/internal/slices"
	"github.com/Tom5521/gotext-tools/v2/internal/util"
	"github.com/Tom5521/gotext-tools/v2/pkg/po"
)

// Aliase this bc I'm too lazy to write "uint32" every time I want to use it.
type u32 = uint32

const (
	eot = "\x04"
	nul = "\x00"
)

func (mc MoCompiler) info(format string, a ...any) {
	if mc.Config.Logger != nil && mc.Config.Verbose {
		mc.Config.Logger.Println("INFO:", fmt.Sprintf(format, a...))
	}
}

func (mc MoCompiler) error(format string, a ...any) error {
	if mc.Config.IgnoreErrors {
		return nil
	}
	format = "compile: " + format
	err := fmt.Errorf(format, a...)
	if mc.Config.Logger != nil {
		mc.Config.Logger.Println("ERROR:", err)
	}

	return err
}

// A len() function with fixed-size return.
func flen(value any) u32 {
	return u32(reflect.ValueOf(value).Len())
}

func max[T slices.Ordered](values ...T) T {
	return slices.Max(values)
}

func (mc *MoCompiler) writeTo(writer io.Writer) error {
	entries := mc.File.Entries.Solve().CleanFuzzy().CleanObsoletes()
	entries = entries.SortFunc(po.CompareEntryByID)

	mc.info("creating header...")
	var hashTabSize u32
	if mc.Config.HashTable {
		hashTabSize = max(3, util.NextPrime((flen(entries)*4)/3))
	}

	header := util.MoHeader{
		Magic:          mc.Config.Endianness.MagicNumber(),
		Revision:       0,
		Nstrings:       flen(entries),
		OrigTabOffset:  7 * 4,
		TransTabOffset: 7*4 + flen(entries)*8,
		HashTabSize:    hashTabSize,
		HashTabOffset:  7*4 + 16*flen(entries),
	}

	mc.info("creating offsets...")
	// Original code translated from: https://github.com/izimobil/polib/blob/master/polib.py#L553
	var (
		offsets   []u32
		ids, strs string
	)

	for _, e := range entries {
		msgid := e.UnifiedID()
		msgstr := e.UnifiedStr()

		offsets = append(offsets,
			flen(ids),
			flen(msgid),
			flen(strs),
			flen(msgstr),
		)
		ids += msgid + nul
		strs += msgstr + nul
	}

	origStart := header.HashTabOffset + (hashTabSize * 4)
	transStart := origStart + flen(ids)

	var origOffsets, transOffsets []u32
	for i := 0; i < len(offsets); i += 4 {
		o1 := offsets[i]
		l1 := offsets[i+1]
		o2 := offsets[i+2]
		l2 := offsets[i+3]

		origOffsets = append(origOffsets, l1, o1+origStart)
		transOffsets = append(transOffsets, l2, o2+transStart)
	}

	mc.info("making hashes...")
	var hashTable []u32
	if mc.Config.HashTable {
		hashTable = buildHashTable(entries, hashTabSize)
	}

	data := []any{
		header,
		origOffsets,
		transOffsets,
		hashTable,
		[]byte(ids),
		[]byte(strs),
	}

	for _, v := range data {
		err := bin.Write(writer, mc.Config.Endianness.Order(), v)
		if err != nil {
			return mc.error("error writing binary data: %w", err)
		}
	}

	return nil
}

// Translated from https://github.com/autotools-mirror/gettext/blob/master/gettext-tools/src/write-mo.c#L876
func buildHashTable(entries po.Entries, size u32) []u32 {
	hashMap := make([]u32, size)

	var seq u32 = 1
	for _, e := range entries {
		hashVal := e.Hash()
		idx := hashVal % size

		if hashMap[idx] != 0 {
			incr := 1 + (hashVal % (size - 2))
			for {
				diff := size - incr
				if idx >= diff {
					idx -= diff
				} else {
					idx += incr
				}
				if hashMap[idx] == 0 {
					break
				}
			}
		}

		hashMap[idx] = seq
		seq++
	}

	return hashMap
}
