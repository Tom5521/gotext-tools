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

	// NOTE: Change the behavior of this, its should be modified
	// to follow the general of the gettext tools.
	//
	// NOTE: (2) It really has to be changed? This library has nothing to do
	// with the gettext tools, it's independent of it. The PoConfig & PoCompiler
	// should be changed to be more general and utilitarian purpose.

	// If true, it still writes to the file if it already exists, in the method [MoCompiler.ToFile].
	Force bool
	// If true, process information and warnings are also printed.
	Verbose      bool
	IgnoreErrors bool
	Endianness   Endianness
	// If true, compiles the hash table.
	HashTable bool

	// NOTE: This reaaaaalyyy need to be exposed?

	// DepureEntries determines whether fuzzy,
	// obsolete, and duplicate entries will be
	// deleted before starting the process.
	//
	// WARNING: Only disable this if you know what you're doing.
	DepureEntries bool
	// WARNING: Only disable this if you know what you're doing.
	SortEntries bool
}

type Endianness = util.Endianness

const (
	LittleEndian = util.LittleEndian
	BigEndian    = util.BigEndian
	NativeEndian = util.NativeEndian
)

// ApplyOptions overwrites the configuration with the options provided,
// saving the previous state so that it can be restored
// later with [MoConfig.RestoreLastCfg] if desired.
func (mc *MoConfig) ApplyOptions(opts ...MoOption) {
	mc.lastCfg = *mc

	for _, opt := range opts {
		opt(mc)
	}
}

// RestoreLastCfg restores the configuration state prior to the last
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
		HashTable:     true,
		DepureEntries: true,
		SortEntries:   true,
	}
	c.ApplyOptions(opts...)
	return c
}

// MoOption defines functions that modify MoConfig.
type MoOption func(c *MoConfig)

// MoWithEndianness sets the byte order for MO file output.
func MoWithEndianness(e Endianness) MoOption {
	return func(c *MoConfig) {
		c.Endianness = e
	}
}

func MoWithDepureEntries(b bool) MoOption {
	return func(c *MoConfig) {
		c.DepureEntries = b
	}
}

func MoWithSortEntries(s bool) MoOption {
	return func(c *MoConfig) {
		c.SortEntries = s
	}
}

// MoWithConfig replaces the entire configuration.
func MoWithConfig(n MoConfig) MoOption {
	return func(c *MoConfig) {
		*c = n
	}
}

// MoWithHashTable toggles hash table generation.
func MoWithHashTable(h bool) MoOption {
	return func(c *MoConfig) {
		c.HashTable = h
	}
}

// MoWithForce toggles file overwrite behavior.
func MoWithForce(f bool) MoOption {
	return func(c *MoConfig) {
		c.Force = f
	}
}

// MoWithIgnoreErrors toggles error suppression.
func MoWithIgnoreErrors(i bool) MoOption {
	return func(c *MoConfig) {
		c.IgnoreErrors = i
	}
}

// MoWithLogger sets the output logger.
func MoWithLogger(l *log.Logger) MoOption {
	return func(c *MoConfig) {
		c.Logger = l
	}
}

// MoWithVerbose toggles detailed logging.
func MoWithVerbose(v bool) MoOption {
	return func(c *MoConfig) {
		c.Verbose = v
	}
}
