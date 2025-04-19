package util

import (
	"encoding/binary"
	"unsafe"
)

const (
	BigEndianMagicNumber    uint32 = 0xde120495
	LittleEndianMagicNumber uint32 = 0x950412de
)

var IsBigEndian = (*[2]uint8)(unsafe.Pointer(&[]uint16{1}[0]))[0] == 0

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
