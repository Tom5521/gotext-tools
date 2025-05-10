package parse

import (
	"bytes"
	bin "encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Tom5521/gotext-tools/v2/internal/util"
	"github.com/Tom5521/gotext-tools/v2/pkg/po"
)

// Type aliases for cleaner code
type (
	u32 = uint32
	i64 = int64
	i32 = int32
)

// Common byte sequences used in MO file parsing
var (
	eot = []byte{4} // End Of Transmission marker
	nul = []byte{0} // Null byte separator
)

// Ensure MoParser implements the po.Parser interface
var _ po.Parser = (*MoParser)(nil)

// MoParser handles parsing of MO files into po.File structures.
type MoParser struct {
	Config MoConfig // Configuration for parsing behavior

	data     []byte  // Raw MO file data
	filename string  // Name of the source file
	errors   []error // Collection of errors encountered during parsing
}

// error logs an error message and adds it to the parser's error collection.
// If a logger is configured in the parser's Config, it will also log the error.
func (m *MoParser) error(format string, a ...any) {
	var err error
	format = "po/parse: " + format
	if len(a) == 0 {
		err = errors.New(format)
	} else {
		err = fmt.Errorf(format, a...)
	}

	if m.Config.Logger != nil {
		m.Config.Logger.Println("ERROR:", err)
	}

	m.errors = append(m.errors, err)
}

// NewMo creates a new MoParser from a file path.
func NewMo(path string, opts ...MoOption) (*MoParser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return NewMoFromReader(f, path, opts...)
}

// NewMoFromReader creates a new MoParser from an io.Reader.
func NewMoFromReader(r io.Reader, name string, opts ...MoOption) (*MoParser, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return &MoParser{data: b, filename: name, Config: DefaultMoConfig(opts...)}, nil
}

// NewMoFromFile creates a new MoParser from an open *os.File.
func NewMoFromFile(f *os.File, opts ...MoOption) (*MoParser, error) {
	return NewMoFromReader(f, f.Name(), opts...)
}

// NewMoFromBytes creates a new MoParser from a byte slice.
func NewMoFromBytes(b []byte, name string, opts ...MoOption) *MoParser {
	return &MoParser{
		data:     b,
		filename: name,
		Config:   DefaultMoConfig(opts...),
	}
}

// Error returns the first error encountered during parsing, if any.
func (m MoParser) Error() error {
	if len(m.errors) == 0 {
		return nil
	}
	return m.errors[0]
}

// Errors returns all errors encountered during parsing.
func (m MoParser) Errors() []error {
	return m.errors
}

// defineOrder determines the byte order (endianness) of the MO file.
// It reads the magic number from the file header to detect the byte order.
func (m *MoParser) defineOrder(reader *bytes.Reader) (order bin.ByteOrder) {
	var magic u32

	if m.Config.Endianness == NativeEndian {
		err := bin.Read(reader, NativeEndian.Order(), &magic)
		if err != nil {
			m.error("error reading magic number: %w", err)
			return
		}
		switch magic {
		case util.LittleEndianMagicNumber:
			order = bin.LittleEndian
		case util.BigEndianMagicNumber:
			order = bin.BigEndian
		default:
			m.error("invalid magic number, this isn't a MO file")
			return
		}
	} else {
		endian := m.Config.Endianness
		order = endian.Order()
		err := bin.Read(reader, order, &magic)
		if err != nil {
			m.error("error reading magic number: %w", err)
			return
		}
		if magic != endian.MagicNumber() {
			m.error("invalid magic number, this isn't a MO file")
		}
	}

	reader.Seek(0, 0)

	return
}

// ParseWithOptions parses the MO file with temporary configuration options.
// The original configuration is restored after parsing.
func (m *MoParser) ParseWithOptions(opts ...MoOption) (file *po.File) {
	m.Config.ApplyOptions(opts...)
	defer m.Config.RestoreLastCfg()

	return m.Parse()
}

// Parse reads and parses the MO file into a po.File structure.
func (m *MoParser) Parse() (file *po.File) {
	var err error
	r := bytes.NewReader(m.data)

	bo := m.defineOrder(r)
	if m.Error() != nil {
		return
	}

	var header util.MoHeader
	if err = bin.Read(r, bo, &header); err != nil {
		m.error("error reading header: %w", err)
		return
	}

	// Validate version numbers
	if v := header.MajorVersion; v != 0 && v != 1 {
		m.error("invalid major version number (%d)", v)
	}

	if v := header.MinorVersion; v != 0 && v != 1 {
		m.error("invalid minor version number (%d)", v)
	}

	// Read message ID table
	msgIDStart := make([]u32, header.Nstrings)
	msgIDLen := make([]u32, header.Nstrings)
	_, err = r.Seek(i64(header.OrigTabOffset), 0)
	if err != nil {
		m.error("bad original table offset(%d): %w", header.OrigTabOffset, err)
		return
	}

	for i := u32(0); i < header.Nstrings; i++ {
		if err = bin.Read(r, bo, &msgIDLen[i]); err != nil {
			m.error("error reading msgid len[%d]: %w", i, err)
			return
		}
		if err = bin.Read(r, bo, &msgIDStart[i]); err != nil {
			m.error("error reading msgid start[%d]: %w", i, err)
			return
		}
	}

	// Read message string table
	msgStrStart := make([]i32, header.Nstrings)
	msgStrLen := make([]i32, header.Nstrings)
	_, err = r.Seek(i64(header.TransTabOffset), 0)
	if err != nil {
		m.error("bad translation table offset(%d): %w", header.TransTabOffset, err)
		return
	}

	for i := u32(0); i < header.Nstrings; i++ {
		if err = bin.Read(r, bo, &msgStrLen[i]); err != nil {
			m.error("error reading msgstr len[%d]: %w", i, err)
			return
		}
		if err = bin.Read(r, bo, &msgStrStart[i]); err != nil {
			m.error("error reading msgstr start[%d]: %w", i, err)
			return
		}
	}

	// Create the po.File structure with all entries
	file = &po.File{
		Name: m.filename,
		Entries: m.makeEntries(
			r,
			&header,
			msgIDStart,
			msgIDLen,
			msgStrStart,
			msgStrLen,
		),
	}

	// Validate sorting if required by configuration
	if m.Config.MustBeSorted {
		if !file.IsSortedFunc(po.CompareEntryByID) {
			m.error("entries must be sorted")
			return
		}
	}

	return
}

// makeEntries creates po.Entries from the parsed MO file data.
func (m *MoParser) makeEntries(
	r *bytes.Reader,
	header *util.MoHeader,
	msgIDStart, msgIDLen []u32,
	msgStrStart, msgStrLen []i32,
) (entries po.Entries) {
	for i := u32(0); i < header.Nstrings; i++ {
		idStart := i64(msgIDStart[i])
		idLen := msgIDLen[i]
		strStart := i64(msgStrStart[i])
		strLen := msgStrLen[i]

		// Read message ID data
		_, err := r.Seek(idStart, 0)
		if err != nil {
			m.error("bad msgid start[%d]: %w", idStart, err)
		}

		msgIDData := make([]byte, idLen)
		_, err = r.Read(msgIDData)
		if err != nil {
			m.error(
				"error reading msgid data[start: %d len: %d]: %w",
				idStart,
				idLen,
				err,
			)
		}

		// Read message string data
		_, err = r.Seek(strStart, 0)
		if err != nil {
			m.error("bad msgstr start[%d]: %w", strStart, err)
		}
		msgStrData := make([]byte, strLen)
		_, err = r.Read(msgStrData)
		if err != nil {
			m.error("error reading msgstr data[start: %d len: %d]: %w",
				strStart,
				strLen,
				err,
			)
		}

		entries = append(entries, makeEntry(msgIDData, msgStrData))
	}

	return
}

// makeEntry creates a single po.Entry from raw message ID and string data.
func makeEntry(msgid, msgstr []byte) (entry po.Entry) {
	var (
		msgctxt  []byte
		pluralID []byte
	)

	// Split context from message ID (separated by EOT)
	d := bytes.Split(msgid, eot)
	if len(d) == 1 {
		msgid = d[0]
	} else {
		msgid, msgctxt = d[1], d[0]
	}

	// Handle plural forms in message ID (separated by NUL)
	msgidParts := bytes.Split(msgid, nul)
	if len(msgidParts) > 1 {
		msgid = msgidParts[0]
		if len(msgidParts) >= 2 {
			pluralID = msgidParts[1]
		}
	}

	entry.ID = string(msgid)
	entry.Plural = string(pluralID)

	// Handle plural forms in message string
	msgstrParts := bytes.Split(msgstr, nul)
	if len(msgstrParts) == 1 {
		entry.Str = string(msgstrParts[0])
	} else {
		for i, s := range msgstrParts {
			entry.Plurals = append(entry.Plurals,
				po.PluralEntry{
					ID:  i,
					Str: string(s),
				},
			)
		}
	}

	// Set context if present
	if len(msgctxt) > 0 {
		entry.Context = string(msgctxt)
	}

	return
}
