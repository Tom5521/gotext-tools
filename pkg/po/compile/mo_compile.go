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

// u32 is an alias for uint32 used for convenience throughout the package.
type u32 = uint32

const (
	// eot represents the End Of Transmission character (ASCII 0x04)
	eot = "\x04"
	// nul represents the null character (ASCII 0x00)
	nul = "\x00"
)

// info logs an informational message if verbose logging is enabled.
// It prepends "INFO:" to the message and sends it to the configured logger.
func (mc MoCompiler) info(format string, a ...any) {
	if mc.Config.Logger != nil && mc.Config.Verbose {
		mc.Config.Logger.Println("INFO:", fmt.Sprintf(format, a...))
	}
}

// error creates and logs an error message. If IgnoreErrors is true, it returns nil.
// The error message is prepended with "compile:" to indicate its origin.
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

// flen returns the length of the given value as a fixed-size u32.
// It uses reflection to get the length of various types.
func flen(value any) u32 {
	return u32(reflect.ValueOf(value).Len())
}

// max returns the maximum value among the provided values of ordered types.
// It uses generics to work with any type that satisfies the Ordered constraint.
func max[T slices.Ordered](values ...T) T {
	return slices.Max(values)
}

// writeTo writes the compiled MO file data to the provided writer.
// It performs several steps:
//  1. Cleans and sorts the PO entries
//  2. Creates the MO file header
//  3. Generates offsets for original and translated strings
//  4. Optionally builds a hash table for faster lookups
//  5. Encodes all data in the specified endianness
//
// Returns an error if any step fails, unless IgnoreErrors is true.
func (mc *MoCompiler) writeTo(writer io.Writer) error {
	mc.info("cleaning & sorting entries...")
	entries := mc.File.Entries.Solve().CleanFuzzy().CleanObsoletes()
	entries = entries.SortFunc(po.CompareEntryByID)

	mc.info("creating header...")
	var hashTabSize u32
	if mc.Config.HashTable {
		hashTabSize = max(3, util.NextPrime((flen(entries)*4)/3))
	}

	header := util.MoHeader{
		Magic:          mc.Config.Endianness.MagicNumber(),
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

	mc.info("encoding...")
	for _, v := range data {
		err := bin.Write(writer, mc.Config.Endianness.Order(), v)
		if err != nil {
			return mc.error("error encoding binary data: %w", err)
		}
	}

	return nil
}

// buildHashTable creates a hash table for the given PO entries using the specified size.
// The implementation is translated from gettext's write-mo.c:
// https://github.com/autotools-mirror/gettext/blob/master/gettext-tools/src/write-mo.c#L876
//
// It uses open addressing with double hashing to handle collisions.
// Each entry's sequence number (1-based) is stored at its calculated position.
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
