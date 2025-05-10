package parse

import (
	"log"

	"github.com/Tom5521/gotext-tools/v2/internal/util"
)

// MoConfig holds configuration options for MO file parsing.
type MoConfig struct {
	// It is used to restore the configuration using the method [MoConfig.RestoreLastCfg]
	// and is saved when using the asd method [MoConfig.ApplyOptions].
	lastCfg any

	// The logger can be nil, otherwise this logger will be used to print all errors by default.
	Logger *log.Logger
	// Specifies the endianness to work with,
	// if [NativeEndian], the endianness
	// automatically determined by the magic number will be used,
	// otherwise, the specified endianness will be used.
	Endianness Endianness
	// Causes parsing to fail if the entries are not sorted properly.
	MustBeSorted bool
}

func DefaultMoConfig(opts ...MoOption) MoConfig {
	mc := MoConfig{}

	mc.ApplyOptions(opts...)
	return mc
}

// Restores the configuration state prior to the last
// [MoConfig.ApplyOptions] if it exists, otherwise it does nothing.
func (mc *MoConfig) RestoreLastCfg() {
	if mc.lastCfg != nil {
		*mc = mc.lastCfg.(MoConfig)
	}
}

// Overwrite the configuration with the options provided,
// saving the previous state so that it can be restored
// later with [MoConfig.RestoreLastCfg] if desired.
func (mc *MoConfig) ApplyOptions(opts ...MoOption) {
	mc.lastCfg = *mc

	for _, mo := range opts {
		mo(mc)
	}
}

// Endianness represents byte order options for MO files.
type Endianness = util.Endianness

const (
	LittleEndian = util.LittleEndian // Force little-endian parsing
	BigEndian    = util.BigEndian    // Force big-endian parsing
	NativeEndian = util.NativeEndian // Auto-detect byte order
)

// MoOption defines a function type for modifying MoConfig.
type MoOption func(*MoConfig)

// MoWithMustBeSorted creates an option to enforce entry sorting.
func MoWithMustBeSorted(m bool) MoOption {
	return func(mc *MoConfig) {
		mc.MustBeSorted = m
	}
}

// MoWithLogger creates an option to set the error logger.
func MoWithLogger(logger *log.Logger) MoOption {
	return func(mc *MoConfig) {
		mc.Logger = logger
	}
}

// MoWithConfig creates an option to replace the entire configuration.
func MoWithConfig(c MoConfig) MoOption {
	return func(mc *MoConfig) { *mc = c }
}

// MoWithEndianness creates an option to set the byte order.
func MoWithEndianness(e Endianness) MoOption {
	return func(mc *MoConfig) { mc.Endianness = e }
}
