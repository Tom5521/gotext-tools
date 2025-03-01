package compiler

import (
	"errors"
	"io"
	"log"
)

// PoConfig holds the settings for the compiler, affecting how translations are processed.
type PoConfig struct {
	Logger          *log.Logger
	ForcePo         bool           // If true, forces the creation of a `.po` file, even if not strictly needed.
	OmitHeader      bool           // If true, omits the header section in the generated `.po` file.
	PackageName     string         // Name of the package associated with the translation.
	CopyrightHolder string         // Name of the entity holding copyright over the translation.
	ForeignUser     bool           // If true, marks the translation as public domain.
	Title           string         // Title to be included in the `.po` file header.
	NoLocation      bool           // If true, suppresses location comments in the `.po` file.
	AddLocation     PoLocationMode // Specifies how location comments should be included ("never", "file", "full").
	MsgstrPrefix    string         // Prefix added to all translation strings.
	MsgstrSuffix    string         // Suffix added to all translation strings.
	IgnoreErrors    bool           // If true, allows compilation to proceed despite non-critical errors.
	Verbose         bool
}

type PoLocationMode string

const (
	PoLocationModeFull  PoLocationMode = "full"
	PoLocationModeNever PoLocationMode = "never"
	PoLocationModeFile  PoLocationMode = "file"
)

func DefaultPoConfig(opts ...PoOption) PoConfig {
	c := PoConfig{
		Logger:      log.New(io.Discard, "", 0),
		PackageName: "PACKAGE NAME",
		AddLocation: PoLocationModeFull,
	}

	for _, opt := range opts {
		opt(&c)
	}

	return c
}

func NewPoConfigFromOptions(opts ...PoOption) PoConfig {
	var config PoConfig

	for _, opt := range opts {
		opt(&config)
	}

	return config
}

type PoOption func(*PoConfig)

func PoWithVerbose(v bool) PoOption {
	return func(c *PoConfig) {
		c.Verbose = v
	}
}

func PoWithLogger(logger *log.Logger) PoOption {
	return func(c *PoConfig) {
		c.Logger = logger
	}
}

func PoWithIgnoreErrors(i bool) PoOption {
	return func(c *PoConfig) {
		c.IgnoreErrors = i
	}
}

func PoWithConfig(cfg PoConfig) PoOption {
	return func(c *PoConfig) {
		*c = cfg
	}
}

func PoWithForcePo(f bool) PoOption {
	return func(c *PoConfig) {
		c.ForcePo = f
	}
}

func PoWithOmitHeader(o bool) PoOption {
	return func(c *PoConfig) {
		c.OmitHeader = o
	}
}

func PoWithPackageName(name string) PoOption {
	return func(c *PoConfig) {
		c.PackageName = name
	}
}

func PoWithCopyrightHolder(holder string) PoOption {
	return func(c *PoConfig) {
		c.CopyrightHolder = holder
	}
}

func PoWithForeignUser(f bool) PoOption {
	return func(c *PoConfig) {
		c.ForeignUser = f
	}
}

func PoWithTitle(t string) PoOption {
	return func(c *PoConfig) {
		c.Title = t
	}
}

func PoWithNoLocation(n bool) PoOption {
	return func(c *PoConfig) {
		c.NoLocation = n
	}
}

func PoWithAddLocation(loc PoLocationMode) PoOption {
	return func(c *PoConfig) {
		c.AddLocation = loc
	}
}

func PoWithMsgstrPrefix(prefix string) PoOption {
	return func(c *PoConfig) {
		c.MsgstrPrefix = prefix
	}
}

func PoWithMsgstrSuffix(suffix string) PoOption {
	return func(c *PoConfig) {
		c.MsgstrSuffix = suffix
	}
}

func (c PoConfig) Validate() error {
	if c.NoLocation && c.AddLocation != PoLocationModeNever {
		return errors.New("noLocation and AddLocation are in conflict")
	}

	return nil
}
