package util

import (
	"encoding/binary"

	"golang.org/x/sys/cpu"
)

const (
	BigEndianMagicNumber    uint32 = 0xde120495
	LittleEndianMagicNumber uint32 = 0x950412de
)

var IsBigEndian = cpu.IsBigEndian

var NativeEndianOrder = func() binary.ByteOrder {
	if IsBigEndian {
		return binary.BigEndian
	}
	return binary.LittleEndian
}()

type Endianness int

const (
	NativeEndian Endianness = iota
	LittleEndian
	BigEndian
)

func (e Endianness) Order() binary.ByteOrder {
	switch e {
	case LittleEndian:
		return binary.LittleEndian
	case BigEndian:
		return binary.BigEndian
	case NativeEndian:
		fallthrough
	default:
		return NativeEndianOrder
	}
}

func (e Endianness) MagicNumber() uint32 {
	switch e {
	case LittleEndian:
		return LittleEndianMagicNumber
	case BigEndian:
		return BigEndianMagicNumber
	case NativeEndian:
		fallthrough
	default:
		if IsBigEndian {
			return BigEndianMagicNumber
		}
		return LittleEndianMagicNumber
	}
}

func (e Endianness) String() string {
	switch e {
	case LittleEndian:
		return "Little"
	case BigEndian:
		return "Big"
	case NativeEndian:
		fallthrough
	default:
		return "Native"
	}
}

type MoHeader struct {
	Magic          u32    // 0
	MajorVersion   uint16 // 2
	MinorVersion   uint16 // 4
	Nstrings       u32    // 8
	OrigTabOffset  u32    // 12
	TransTabOffset u32    // 16
	HashTabSize    u32    // 20
	HashTabOffset  u32    // 24
}
