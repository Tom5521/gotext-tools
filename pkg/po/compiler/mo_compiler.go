//go:build experimental
// +build experimental

package compiler

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	"github.com/Tom5521/xgotext/pkg/po"
)

// var _ Compiler = (*MoCompiler)(nil)

var (
	isBigEndian = (*(*[2]uint8)(unsafe.Pointer(&[]uint16{1}[0])))[0] == 0

	magicNumber = func() uint32 {
		if isBigEndian {
			return bigEndianMagicNumber
		}
		return littleEndianMagicNumber
	}()
	endian = binary.NativeEndian
)

const (
	bigEndianMagicNumber    uint32 = 0xde120495
	littleEndianMagicNumber uint32 = 0x950412de

	revVersionMajor, revVersionMinor uint16 = 0x0001, 0x0000
)

func makeMagicNumber() []byte {
	b := make([]byte, 4)
	endian.AppendUint32(b, magicNumber)

	return b
}

func makeRevVersions() []byte {
	var buf []byte
	a, b := make([]byte, 2), make([]byte, 2)
	endian.AppendUint16(a, revVersionMajor)
	endian.AppendUint16(b, revVersionMinor)

	buf = append(buf, a...)
	buf = append(buf, b...)

	return buf
}

type MoCompiler struct {
	File *po.File
}

func (mc *MoCompiler) applyOptions(opts ...Option) {
	for _, opt := range opts {
		opt(&mc.Config)
	}
}

func (mc *MoCompiler) createBinary() {
	var buf bytes.Buffer

	buf.Write(makeMagicNumber())
	buf.Write(makeRevVersions())
}
