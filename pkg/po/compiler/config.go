package compiler

import "errors"

// Config holds the settings for the compiler, affecting how translations are processed.
type Config struct {
	ForcePo         bool         // If true, forces the creation of a `.po` file, even if not strictly needed.
	OmitHeader      bool         // If true, omits the header section in the generated `.po` file.
	PackageName     string       // Name of the package associated with the translation.
	CopyrightHolder string       // Name of the entity holding copyright over the translation.
	ForeignUser     bool         // If true, marks the translation as public domain.
	Title           string       // Title to be included in the `.po` file header.
	NoLocation      bool         // If true, suppresses location comments in the `.po` file.
	AddLocation     LocationMode // Specifies how location comments should be included ("never", "file", "full").
	MsgstrPrefix    string       // Prefix added to all translation strings.
	MsgstrSuffix    string       // Suffix added to all translation strings.
	IgnoreErrors    bool         // If true, allows compilation to proceed despite non-critical errors.
}

type LocationMode string

const (
	LocationModeFull  LocationMode = "full"
	LocationModeNever LocationMode = "never"
	LocationModeFile  LocationMode = "file"
)

func NewConfigFromOptions(opts ...Option) Config {
	var config Config

	for _, opt := range opts {
		opt(&config)
	}

	return config
}

type Option func(*Config)

func WithIgnoreErrors(i bool) Option {
	return func(c *Config) {
		c.IgnoreErrors = i
	}
}

func WithConfig(cfg Config) Option {
	return func(c *Config) {
		*c = cfg
	}
}

func WithForcePo(f bool) Option {
	return func(c *Config) {
		c.ForcePo = f
	}
}

func WithOmitHeader(o bool) Option {
	return func(c *Config) {
		c.OmitHeader = o
	}
}

func WithPackageName(name string) Option {
	return func(c *Config) {
		c.PackageName = name
	}
}

func WithCopyrightHolder(holder string) Option {
	return func(c *Config) {
		c.CopyrightHolder = holder
	}
}

func WithForeignUser(f bool) Option {
	return func(c *Config) {
		c.ForeignUser = f
	}
}

func WithTitle(t string) Option {
	return func(c *Config) {
		c.Title = t
	}
}

func WithNoLocation(n bool) Option {
	return func(c *Config) {
		c.NoLocation = n
	}
}

func WithAddLocation(loc LocationMode) Option {
	return func(c *Config) {
		c.AddLocation = loc
	}
}

func WithMsgstrPrefix(prefix string) Option {
	return func(c *Config) {
		c.MsgstrPrefix = prefix
	}
}

func WithMsgstrSuffix(suffix string) Option {
	return func(c *Config) {
		c.MsgstrSuffix = suffix
	}
}

func (c Config) Validate() error {
	if c.NoLocation && c.AddLocation != LocationModeNever {
		return errors.New("noLocation and AddLocation are in conflict")
	}

	return nil
}
