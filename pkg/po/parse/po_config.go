package parse

import (
	"log"
)

type PoConfig struct {
	// It is used to restore the configuration using the method [PoConfig.RestoreLastCfg]
	// and is saved when using the asd method [PoConfig.ApplyOptions].
	lastCfg any

	IgnoreComments    bool
	IgnoreAllComments bool
	// The logger can be nil, otherwise this logger will be used to print all errors by default.
	Logger          *log.Logger
	Verbose         bool
	SkipHeader      bool
	CleanDuplicates bool

	ParseObsoletes          bool
	UseCustomObsoletePrefix bool
	CustomObsoletePrefix    rune

	markAllAsObsolete bool
}

// Restores the configuration state prior to the last
// [PoConfig.ApplyOptions] if it exists, otherwise it does nothing.
func (p *PoConfig) RestoreLastCfg() {
	if p.lastCfg != nil {
		*p = p.lastCfg.(PoConfig)
	}
}

// Overwrite the configuration with the options provided,
// saving the previous state so that it can be restored
// later with [PoConfig.RestoreLastCfg] if desired.
func (p *PoConfig) ApplyOptions(opts ...PoOption) {
	p.lastCfg = *p
	for _, opt := range opts {
		opt(p)
	}
}

func DefaultPoConfig(opts ...PoOption) PoConfig {
	c := PoConfig{
		CleanDuplicates: true,
		ParseObsoletes:  true,
	}

	c.ApplyOptions(opts...)

	return c
}

type PoOption func(*PoConfig)

func poWithMarkAllAsObsolete(m bool) PoOption {
	return func(pc *PoConfig) {
		pc.markAllAsObsolete = m
	}
}

func PoWithVerbose(v bool) PoOption {
	return func(pc *PoConfig) {
		pc.Verbose = v
	}
}

func PoWithParseObsolete(p bool) PoOption {
	return func(pc *PoConfig) {
		pc.ParseObsoletes = p
	}
}

func PoWithCustomObsoletePrefix(r rune) PoOption {
	return func(pc *PoConfig) {
		pc.CustomObsoletePrefix = r
	}
}

func PoWithUseCustomObsoletePrefix(u bool) PoOption {
	return func(pc *PoConfig) {
		pc.UseCustomObsoletePrefix = u
	}
}

func PoWithIgnoreAllComments(iag bool) PoOption {
	return func(pc *PoConfig) {
		pc.IgnoreAllComments = iag
	}
}

func PoWithIgnoreComments(ig bool) PoOption {
	return func(pc *PoConfig) {
		pc.IgnoreComments = ig
	}
}

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
