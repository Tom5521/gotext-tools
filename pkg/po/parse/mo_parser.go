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

type (
	u32 = uint32
	i64 = int64
	i32 = int32
)

var (
	eot = []byte{4}
	nul = []byte{0}
)

var _ po.Parser = (*MoParser)(nil)

type MoParser struct {
	data     []byte
	filename string
	errors   []error
	Config   MoConfig
}

// This method MUST be used to log any errors inside this structure.
func (m *MoParser) error(format string, a ...any) {
	var err error
	format = "parse: " + format
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

func NewMo(path string, opts ...MoOption) (*MoParser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return NewMoFromReader(f, path, opts...)
}

func NewMoFromReader(r io.Reader, name string, opts ...MoOption) (*MoParser, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return &MoParser{data: b, filename: name, Config: DefaultMoConfig(opts...)}, nil
}

func NewMoFromFile(f *os.File, opts ...MoOption) (*MoParser, error) {
	return NewMoFromReader(f, f.Name(), opts...)
}

func NewMoFromBytes(b []byte, name string, opts ...MoOption) *MoParser {
	return &MoParser{
		data:     b,
		filename: name,
		Config:   DefaultMoConfig(opts...),
	}
}

// Return the first error in the stack.
func (m MoParser) Error() error {
	if len(m.errors) == 0 {
		return nil
	}

	return m.errors[0]
}

func (m MoParser) Errors() []error {
	return m.errors
}

func (m *MoParser) genBasics(reader *bytes.Reader) (order bin.ByteOrder, err error) {
	var magicNumber u32
	if m.Config.Endianness == NativeEndian {
		if err = bin.Read(reader, bin.LittleEndian, &magicNumber); err != nil {
			m.error("error reading magic number: %w", err)
			return
		}
		switch magicNumber {
		case util.LittleEndianMagicNumber:
			order = bin.LittleEndian
		case util.BigEndianMagicNumber:
			order = bin.BigEndian
		default:
			m.error("invalid magic number")
			return
		}
	} else {
		order = m.Config.Endianness.Order()
		if err = bin.Read(reader, order, &magicNumber); err != nil {
			m.error("error reading magic number: %w", err)
			return
		}
		if magicNumber != m.Config.Endianness.MagicNumber() {
			m.error("invalid magic number")
		}
	}

	reader.Seek(0, 0)

	return
}

func (m *MoParser) ParseWithOptions(opts ...MoOption) (file *po.File) {
	m.Config.ApplyOptions(opts...)
	defer m.Config.RestoreLastCfg()

	return m.Parse()
}

func (m *MoParser) Parse() (file *po.File) {
	r := bytes.NewReader(m.data)

	bo, err := m.genBasics(r)
	if err != nil {
		return
	}

	var header util.MoHeader
	if err = bin.Read(r, bo, &header); err != nil {
		m.error("error reading header: %w", err)
		return
	}

	if v := header.Revision >> 16; v != 0 && v != 1 {
		m.error("invalid major version number (%d)", v)
	}

	if v := header.Revision & 0xFFFF; v != 0 && v != 1 {
		m.error("invalid minor version number (%d)", v)
	}

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

	if m.Config.MustBeSorted {
		if !file.IsSortedFunc(po.CompareEntryByID) {
			m.error("entries must be sorted")
			return
		}
	}

	return
}

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

func makeEntry(msgid, msgstr []byte) (entry po.Entry) {
	var (
		msgctxt  []byte
		pluralID []byte
	)

	d := bytes.Split(msgid, eot)
	if len(d) == 1 {
		msgid = d[0]
	} else {
		msgid, msgctxt = d[1], d[0]
	}

	msgidParts := bytes.Split(msgid, nul)
	if len(msgidParts) > 1 {
		msgid = msgidParts[0]
		if len(msgidParts) >= 2 {
			pluralID = msgidParts[1]
		}
	}

	entry.ID = string(msgid)
	entry.Plural = string(pluralID)

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

	if len(msgctxt) > 0 {
		entry.Context = string(msgctxt)
	}

	return
}
