package parse

import (
	"log"
)

type PoConfig struct {
	// NOTE: Two copies of the same struct ALWAYS? Good.
	// They REAALYY need a backup?

	// It is used to restore the configuration using the method [PoConfig.RestoreLastCfg]
	// and is saved when using the asd method [PoConfig.ApplyOptions].
	lastCfg any

	// IgnoreComments controls whether to discard translator comments.
	IgnoreComments bool
	// IgnoreAllComments controls whether to discard all comments including extracted comments and flags.
	IgnoreAllComments bool
	// Logger is an optional logger for error output. If nil, errors are only stored internally.
	Logger *log.Logger
	// Verbose enables more detailed logging when true.
	Verbose bool
	// SkipHeader controls whether to skip the metadata header entry.
	SkipHeader bool

	// TODO: Remove this because it's just useless and unnecessary.

	// CleanDuplicates controls whether to remove duplicate entries during parsing.
	CleanDuplicates bool
	// ParseObsoletes controls whether to parse obsolete entries (marked with #~).
	ParseObsoletes bool
	// UseCustomObsoletePrefix enables using a custom prefix for obsolete entries.
	UseCustomObsoletePrefix bool
	// CustomObsoletePrefix defines the custom marker for obsolete entries.
	CustomObsoletePrefix rune
	// markAllAsObsolete internally marks all parsed entries as obsolete.
	markAllAsObsolete bool
}

// RestoreLastCfg restores the configuration state prior to the last
// [PoConfig.ApplyOptions] if it exists, otherwise it does nothing.
func (p *PoConfig) RestoreLastCfg() {
	if p.lastCfg != nil {
		*p = p.lastCfg.(PoConfig)
	}
}

// ApplyOptions overwrite the configuration with the options provided,
// saving the previous state so that it can be restored
// later with [PoConfig.RestoreLastCfg] if desired.
func (p *PoConfig) ApplyOptions(opts ...PoOption) {
	p.lastCfg = *p
	for _, opt := range opts {
		opt(p)
	}
}

// DefaultPoConfig returns a new PoConfig with recommended defaults:
// - CleanDuplicates: true
// - ParseObsoletes: true.
func DefaultPoConfig(opts ...PoOption) PoConfig {
	c := PoConfig{
		CleanDuplicates: true,
		ParseObsoletes:  true,
	}
	c.ApplyOptions(opts...)
	return c
}

// PoOption defines a function type for modifying PoConfig.
type PoOption func(*PoConfig)

// poWithMarkAllAsObsolete creates an option to mark all entries as obsolete.
func poWithMarkAllAsObsolete(m bool) PoOption {
	return func(pc *PoConfig) {
		pc.markAllAsObsolete = m
	}
}

// PoWithVerbose creates an option to enable verbose logging.
func PoWithVerbose(v bool) PoOption {
	return func(pc *PoConfig) {
		pc.Verbose = v
	}
}

// PoWithParseObsolete creates an option to control obsolete entry parsing.
func PoWithParseObsolete(p bool) PoOption {
	return func(pc *PoConfig) {
		pc.ParseObsoletes = p
	}
}

// PoWithCustomObsoletePrefix creates an option to set a custom obsolete marker.
func PoWithCustomObsoletePrefix(r rune) PoOption {
	return func(pc *PoConfig) {
		pc.CustomObsoletePrefix = r
	}
}

// PoWithUseCustomObsoletePrefix creates an option to enable custom obsolete markers.
func PoWithUseCustomObsoletePrefix(u bool) PoOption {
	return func(pc *PoConfig) {
		pc.UseCustomObsoletePrefix = u
	}
}

// PoWithIgnoreAllComments creates an option to discard all comments.
func PoWithIgnoreAllComments(iag bool) PoOption {
	return func(pc *PoConfig) {
		pc.IgnoreAllComments = iag
	}
}

// PoWithIgnoreComments creates an option to discard translator comments.
func PoWithIgnoreComments(ig bool) PoOption {
	return func(pc *PoConfig) {
		pc.IgnoreComments = ig
	}
}

// PoWithConfig creates an option to replace the entire configuration.
func PoWithConfig(cfg PoConfig) PoOption {
	return func(c *PoConfig) { *c = cfg }
}

// PoWithSkipHeader creates an option to skip the metadata header.
func PoWithSkipHeader(s bool) PoOption {
	return func(c *PoConfig) { c.SkipHeader = s }
}

// PoWithCleanDuplicates creates an option to control duplicate removal.
func PoWithCleanDuplicates(cd bool) PoOption {
	return func(c *PoConfig) { c.CleanDuplicates = cd }
}

// PoWithLogger creates an option to set the error logger.
func PoWithLogger(logger *log.Logger) PoOption {
	return func(c *PoConfig) { c.Logger = logger }
}
