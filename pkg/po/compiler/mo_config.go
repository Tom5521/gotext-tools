package compiler

import (
	"io"
	"log"
)

type MoConfig struct {
	Logger       *log.Logger
	Verbose      bool
	IgnoreErrors bool
}

func DefaultMoConfig(opts ...MoOption) MoConfig {
	c := MoConfig{
		Logger: log.New(io.Discard, "", 0),
	}

	for _, opt := range opts {
		opt(&c)
	}

	return c
}

type MoOption func(c *MoConfig)

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
