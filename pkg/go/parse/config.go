package parse

import "github.com/Tom5521/xgotext/pkg/po/types"

type Config struct {
	Exclude       []string
	ExtractAll    bool
	HeaderConfig  *types.HeaderConfig
	HeaderOptions []types.HeaderOption
	Header        *types.Header
	FuzzyMatch    bool
}

func defaultConfig() Config {
	return Config{
		Header: func() *types.Header {
			h := types.DefaultHeader()
			return &h
		}(),
	}
}

type Option func(c *Config)

func WithConfig(cfg Config) Option {
	return func(c *Config) {
		*c = cfg
	}
}

func WithExclude(exclude []string) Option {
	return func(c *Config) {
		c.Exclude = exclude
	}
}

func WithExtractAll(e bool) Option {
	return func(c *Config) {
		c.ExtractAll = e
	}
}

func WithHeaderConfig(h *types.HeaderConfig) Option {
	return func(c *Config) {
		c.HeaderConfig = h
	}
}

func WithHeaderOptions(hopts ...types.HeaderOption) Option {
	return func(c *Config) {
		c.HeaderOptions = hopts
	}
}

func WithHeader(h *types.Header) Option {
	return func(c *Config) {
		c.Header = h
	}
}

func WithFuzzyMatch(f bool) Option {
	return func(c *Config) {
		c.FuzzyMatch = f
	}
}
