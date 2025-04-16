package compile

import (
	"io"
	"log"

	"github.com/Tom5521/gotext-tools/v2/pkg/po"
)

type MoConfig struct {
	lastCfg any

	Logger       *log.Logger
	Force        bool
	Verbose      bool
	IgnoreErrors bool
	Sort         bool
	SortMode     po.SortMode
	// HashTable    bool
}

func (mc *MoConfig) ApplyOptions(opts ...MoOption) {
	mc.lastCfg = *mc

	for _, opt := range opts {
		opt(mc)
	}
}

func (mc *MoConfig) RestoreLastCfg() {
	if mc.lastCfg != nil {
		*mc = mc.lastCfg.(MoConfig)
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

func MoWithSortMode(m po.SortMode) MoOption {
	return func(c *MoConfig) {
		c.SortMode = m
	}
}

func MoWithSort(s bool) MoOption {
	return func(c *MoConfig) {
		c.Sort = s
	}
}

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
