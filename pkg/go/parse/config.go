package parse

import (
	"io"
	"log"

	"github.com/Tom5521/gotext-tools/v2/pkg/po"
)

type Config struct {
	// It is used to restore the configuration using the method [Config.RestoreLastCfg]
	// and is saved when using the asd method [Config.ApplyOptions]
	lastCfg any

	Exclude         []string
	ExtractAll      bool
	NoHeader        bool
	HeaderConfig    *po.HeaderConfig
	CustomHeader    *po.Header
	Logger          *log.Logger
	Verbose         bool
	CleanDuplicates bool
}

// Restores the configuration state prior to the last
// [Config.ApplyOptions] if it exists, otherwise it does nothing.
func (c *Config) RestoreLastCfg() {
	if c.lastCfg != nil {
		*c = c.lastCfg.(Config)
	}
}

// Overwrite the configuration with the options provided,
// saving the previous state so that it can be restored
// later with [Config.RestoreLastCfg] if desired.
func (c *Config) ApplyOptions(opts ...Option) {
	c.lastCfg = *c

	for _, opt := range opts {
		opt(c)
	}
}

func DefaultConfig(opts ...Option) Config {
	c := Config{
		Logger:          log.New(io.Discard, "", 0),
		CleanDuplicates: true,
	}

	for _, opt := range opts {
		opt(&c)
	}

	return c
}

type Option func(c *Config)

func WithNoHeader(h bool) Option {
	return func(c *Config) {
		c.NoHeader = h
	}
}

func WithVerbose(v bool) Option {
	return func(c *Config) { c.Verbose = v }
}

func WithLogger(l *log.Logger) Option {
	return func(c *Config) { c.Logger = l }
}

func WithConfig(cfg Config) Option {
	return func(c *Config) { *c = cfg }
}

func WithCleanDuplicates(cl bool) Option {
	return func(c *Config) { c.CleanDuplicates = cl }
}

func WithExclude(exclude ...string) Option {
	return func(c *Config) { c.Exclude = exclude }
}

func WithExtractAll(e bool) Option {
	return func(c *Config) { c.ExtractAll = e }
}

func WithCustomHeader(h *po.Header) Option {
	return func(c *Config) { c.CustomHeader = h }
}

func WithHeaderConfig(h *po.HeaderConfig) Option {
	return func(c *Config) { c.HeaderConfig = h }
}
