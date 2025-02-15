package parse

import (
	"io"
	"log"
)

type Config struct {
	Logger          *log.Logger
	Verbose         bool
	SkipHeader      bool
	CleanDuplicates bool
}

func DefaultConfig(opts ...Option) Config {
	c := Config{
		Logger:          log.New(io.Discard, "", 0),
		CleanDuplicates: true,
	}

	for _, o := range opts {
		o(&c)
	}

	return c
}

type Option func(*Config)

func WithConfig(cfg Config) Option {
	return func(c *Config) { *c = cfg }
}

func WithVerbose(v bool) Option {
	return func(c *Config) { c.Verbose = v }
}

func WithSkipHeader(s bool) Option {
	return func(c *Config) { c.SkipHeader = s }
}

func WithCleanDuplicates(cd bool) Option {
	return func(c *Config) { c.CleanDuplicates = cd }
}

func WithLogger(logger *log.Logger) Option {
	return func(c *Config) { c.Logger = logger }
}
