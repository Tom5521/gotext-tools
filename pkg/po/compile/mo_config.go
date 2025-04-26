package compile

import (
	"io"
	"log"

	"github.com/Tom5521/gotext-tools/v2/internal/util"
)

type MoConfig struct {
	lastCfg any

	Logger       *log.Logger
	Force        bool
	Verbose      bool
	IgnoreErrors bool
	Endianness   Endianness
	HashTable    bool
}

type Endianness = util.Endianness

const (
	LittleEndian = util.LittleEndian
	BigEndian    = util.BigEndian
	NativeEndian = util.NativeEndian
)

func (mc *MoConfig) ApplyOptions(opts ...MoOption) {
	mc.lastCfg = *mc

	for _, opt := range opts {
		opt(mc)
	}
}

func (mc *MoConfig) RestoreLastCfg() {
	if mc.lastCfg != nil {
		*mc = mc.lastCfg.(MoConfig)
	}
}

func DefaultMoConfig(opts ...MoOption) MoConfig {
	c := MoConfig{
		Logger:     log.New(io.Discard, "", 0),
		Endianness: NativeEndian,
	}

	c.ApplyOptions(opts...)

	return c
}

type MoOption func(c *MoConfig)

func MoWithEndianness(e Endianness) MoOption {
	return func(c *MoConfig) {
		c.Endianness = e
	}
}

func MoWithConfig(n MoConfig) MoOption {
	return func(c *MoConfig) {
		*c = n
	}
}

func MoWithHashTable(h bool) MoOption {
	return func(c *MoConfig) {
		c.HashTable = h
	}
}

func MoWithForce(f bool) MoOption {
	return func(c *MoConfig) {
		c.Force = f
	}
}

func MoWithIgnoreErrors(i bool) MoOption {
	return func(c *MoConfig) {
		c.IgnoreErrors = i
	}
}

func MoWithLogger(l *log.Logger) MoOption {
	return func(c *MoConfig) {
		c.Logger = l
	}
}

func MoWithVerbose(v bool) MoOption {
	return func(c *MoConfig) {
		c.Verbose = v
	}
}
