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

var NativeEndian = func() binary.ByteOrder {
	if IsBigEndian {
		return binary.BigEndian
	}
	return binary.LittleEndian
}()
