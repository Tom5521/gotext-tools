package compile

import (
	"log"

	"github.com/Tom5521/gotext-tools/v2/internal/util"
)

type MoConfig struct {
	// It is used to restore the configuration using the method [MoConfig.RestoreLastCfg]
	// and is saved when using the asd method [MoConfig.ApplyOptions]
	lastCfg any

	// The logger can be nil, otherwise this logger will be used to print all errors by default.
	Logger *log.Logger
	// If true, it still writes to the file if it already exists, in the method [MoCompiler.ToFile].
	Force bool
	// If true, process information and warnings are also printed.
	Verbose      bool
	IgnoreErrors bool
	Endianness   Endianness
	// If true, compiles the hash table.
	HashTable bool
}

type Endianness = util.Endianness

const (
	LittleEndian = util.LittleEndian
	BigEndian    = util.BigEndian
	NativeEndian = util.NativeEndian
)

// Overwrite the configuration with the options provided,
// saving the previous state so that it can be restored
// later with [MoConfig.RestoreLastCfg] if desired.
func (mc *MoConfig) ApplyOptions(opts ...MoOption) {
	mc.lastCfg = *mc

	for _, opt := range opts {
		opt(mc)
	}
}

// Restores the configuration state prior to the last
// [MoConfig.ApplyOptions] if it exists, otherwise it does nothing.
func (mc *MoConfig) RestoreLastCfg() {
	if mc.lastCfg != nil {
		*mc = mc.lastCfg.(MoConfig)
	}
}

// DefaultMoConfig creates a new MoConfig with default values.
// Applies any provided options during creation.
func DefaultMoConfig(opts ...MoOption) MoConfig {
	c := MoConfig{
		HashTable: true,
	}
	c.ApplyOptions(opts...)
	return c
}

// MoOption defines functions that modify MoConfig
type MoOption func(c *MoConfig)

// MoWithEndianness sets the byte order for MO file output
func MoWithEndianness(e Endianness) MoOption {
	return func(c *MoConfig) {
		c.Endianness = e
	}
}

// MoWithConfig replaces the entire configuration
func MoWithConfig(n MoConfig) MoOption {
	return func(c *MoConfig) {
		*c = n
	}
}

// MoWithHashTable toggles hash table generation
func MoWithHashTable(h bool) MoOption {
	return func(c *MoConfig) {
		c.HashTable = h
	}
}

// MoWithForce toggles file overwrite behavior
func MoWithForce(f bool) MoOption {
	return func(c *MoConfig) {
		c.Force = f
	}
}

// MoWithIgnoreErrors toggles error suppression
func MoWithIgnoreErrors(i bool) MoOption {
	return func(c *MoConfig) {
		c.IgnoreErrors = i
	}
}

// MoWithLogger sets the output logger
func MoWithLogger(l *log.Logger) MoOption {
	return func(c *MoConfig) {
		c.Logger = l
	}
}

// MoWithVerbose toggles detailed logging
func MoWithVerbose(v bool) MoOption {
	return func(c *MoConfig) {
		c.Verbose = v
	}
}
