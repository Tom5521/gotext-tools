package parse

import (
	"bytes"
	bin "encoding/binary"
	"errors"
	"io"
	"os"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po"
)

type (
	u32 = uint32
	u16 = uint16
	i64 = int64
	i32 = int32
)

var (
	eot = []byte{4}
	nul = []byte{0}
)

type MoParser struct {
	data     io.Reader
	filename string
	errors   []error
}

func NewMo(path string) (*MoParser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return NewMoFromReader(f, path)
}

func NewMoFromReader(r io.Reader, name string) (*MoParser, error) {
	return &MoParser{data: r, filename: name}, nil
}

func NewMoFromFile(f *os.File) (*MoParser, error) {
	return NewMoFromReader(f, f.Name())
}

func NewMoFromBytes(b []byte, name string) (*MoParser, error) {
	return &MoParser{
		data:     bytes.NewReader(b),
		filename: name,
	}, nil
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

func (m *MoParser) genBasics() (reader io.ReadSeeker, order bin.ByteOrder, err error) {
	if b, ok := m.data.(io.ReadSeeker); ok {
		reader = b
	} else {
		var bytedata []byte
		bytedata, err = io.ReadAll(m.data)
		if err != nil {
			return nil, nil, err
		}

		reader = bytes.NewReader(bytedata)
	}

	var magicNumber uint32
	if err = bin.Read(reader, bin.LittleEndian, &magicNumber); err != nil {
		return nil, nil, err
	}
	switch magicNumber {
	case util.LittleEndianMagicNumber:
		order = bin.LittleEndian
	case util.BigEndianMagicNumber:
		order = bin.BigEndian
	default:
		m.errors = append(m.errors, errors.New("invalid magic number"))
		return
	}

	return
}

type moHeader struct {
	MajorVersion u16
	MinorVersion u16
	MsgIDCount   u32
	MsgIDOffset  u32
	MsgStrOffset u32
	HashSize     u32
	HashOffset   u32
}

func (m *MoParser) Parse() (file *po.File) {
	r, bo, err := m.genBasics()
	if err != nil {
		m.errors = append(m.errors, err)
		return
	}

	var header moHeader
	if err = bin.Read(r, bo, &header); err != nil {
		m.errors = append(m.errors, err)
		return
	}

	if v := header.MajorVersion; v != 0 && v != 1 {
		m.errors = append(m.errors, errors.New("invalid version number"))
	}

	if v := header.MinorVersion; v != 0 && v != 1 {
		m.errors = append(m.errors, errors.New("invalid version number"))
	}

	msgIDStart := make([]u32, header.MsgIDCount)
	msgIDLen := make([]u32, header.MsgIDCount)
	if _, err = r.Seek(i64(header.MsgIDOffset), 0); err != nil {
		m.errors = append(m.errors, err)
		return
	}

	for i := u32(0); i < header.MsgIDCount; i++ {
		if err = bin.Read(r, bo, &msgIDLen[i]); err != nil {
			m.errors = append(m.errors, err)
			return
		}
		if err = bin.Read(r, bo, &msgIDStart[i]); err != nil {
			m.errors = append(m.errors, err)
			return
		}
	}

	msgStrStart := make([]i32, header.MsgIDCount)
	msgStrLen := make([]i32, header.MsgIDCount)
	if _, err = r.Seek(i64(header.MsgStrOffset), 0); err != nil {
		m.errors = append(m.errors, err)
		return
	}

	for i := u32(0); i < header.MsgIDCount; i++ {
		if err = bin.Read(r, bo, &msgStrLen[i]); err != nil {
			m.errors = append(m.errors, err)
			return
		}
		if err = bin.Read(r, bo, &msgStrStart[i]); err != nil {
			m.errors = append(m.errors, err)
			return
		}
	}

	file = &po.File{
		Name: m.filename,
		Entries: m.makeEntries(
			r,
			header,
			msgIDStart,
			msgIDLen,
			msgStrStart,
			msgStrLen,
		),
	}

	return
}

func (m *MoParser) makeEntries(
	r io.ReadSeeker,
	header moHeader,
	msgIDStart, msgIDLen []u32,
	msgStrStart, msgStrLen []i32,
) (entries po.Entries) {
	var err error
	for i := u32(0); i < header.MsgIDCount; i++ {
		if _, err = r.Seek(i64(msgIDStart[i]), 0); err != nil {
			m.errors = append(m.errors, err)
			return
		}
		msgIDData := make([]byte, msgIDLen[i])
		if _, err = r.Read(msgIDData); err != nil {
			m.errors = append(m.errors, err)
			return
		}

		if _, err = r.Seek(i64(msgStrStart[i]), 0); err != nil {
			m.errors = append(m.errors, err)
			return
		}

		msgStrData := make([]byte, msgStrLen[i])
		if _, err = r.Read(msgStrData); err != nil {
			m.errors = append(m.errors, err)
			return
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
