package parse

import (
	"io"
	"log"
)

type PoConfig struct {
	Logger          *log.Logger
	SkipHeader      bool
	CleanDuplicates bool
}

func (p *PoConfig) ApplyOptions(opts ...PoOption) {
	for _, opt := range opts {
		opt(p)
	}
}

func DefaultPoConfig(opts ...PoOption) PoConfig {
	c := PoConfig{
		Logger:          log.New(io.Discard, "", 0),
		CleanDuplicates: true,
	}

	for _, o := range opts {
		o(&c)
	}

	return c
}

type PoOption func(*PoConfig)

func PoWithConfig(cfg PoConfig) PoOption {
	return func(c *PoConfig) { *c = cfg }
}

func PoWithSkipHeader(s bool) PoOption {
	return func(c *PoConfig) { c.SkipHeader = s }
}

func PoWithCleanDuplicates(cd bool) PoOption {
	return func(c *PoConfig) { c.CleanDuplicates = cd }
}

func PoWithLogger(logger *log.Logger) PoOption {
	return func(c *PoConfig) { c.Logger = logger }
}
