package parse

import "github.com/Tom5521/gotext-tools/v2/internal/util"

type MoConfig struct {
	lastCfg    any
	Endianness Endianness
}

func DefaultMoConfig(opts ...MoOption) MoConfig {
	mc := MoConfig{
		Endianness: NativeEndian,
	}

	mc.ApplyOptions(opts...)
	return mc
}

func (mc *MoConfig) RestoreLastCfg() {
	if mc.lastCfg != nil {
		*mc = mc.lastCfg.(MoConfig)
	}
}

func (mc *MoConfig) ApplyOptions(opts ...MoOption) {
	mc.lastCfg = *mc

	for _, mo := range opts {
		mo(mc)
	}
}

type Endianness = util.Endianness

const (
	LittleEndian = util.LittleEndian
	BigEndian    = util.BigEndian
	NativeEndian = util.NativeEndian
)

type MoOption func(*MoConfig)

func MoWithConfig(c MoConfig) MoOption {
	return func(mc *MoConfig) { *mc = c }
}

func MoWithEndianness(e Endianness) MoOption {
	return func(mc *MoConfig) { mc.Endianness = e }
}
