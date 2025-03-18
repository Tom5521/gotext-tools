package compiler

import (
	"io"
	"log"
)

type MoConfig struct {
	Logger       *log.Logger
	Force        bool
	Verbose      bool
	IgnoreErrors bool
	// HashTable    bool // TODO: Implement this.
}

func (mc *MoConfig) ApplyOptions(opts ...MoOption) {
	for _, opt := range opts {
		opt(mc)
	}
}

func DefaultMoConfig(opts ...MoOption) MoConfig {
	c := MoConfig{
		Logger: log.New(io.Discard, "", 0),
	}

	c.ApplyOptions(opts...)

	return c
}

type MoOption func(c *MoConfig)

func MoWithConfig(n MoConfig) MoOption {
	return func(c *MoConfig) {
		*c = n
	}
}

// TODO: Finish this.
// func MoWithHashTable(h bool) MoOption {
// 	return func(c *MoConfig) {
// 		c.HashTable = h
// 	}
// }

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
