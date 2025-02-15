package parse

import (
	"io"
	"log"

	"github.com/Tom5521/xgotext/pkg/po/types"
)

// Config defines the configuration options for customizing the parsing process.
//
// ### Attributes:
// - `Exclude`: A list of file paths or directories to exclude from processing.
// - `ExtractAll`: If true, extracts all string literals, not just those marked for translation.
// - `HeaderConfig`: Configuration for the PO file header.
// - `HeaderOptions`: Additional options to customize the PO file header.
// - `Header`: The PO file header to include in the output.
// - `FuzzyMatch`: Enables fuzzy matching for translation entries (e.g., for deduplication).
// - `Logger`: A logger instance for tracking parsing activity and errors.
// - `Verbose`: If true, enables verbose logging.
//
// ### Responsibilities:
// - Provide flexible options to customize the parser behavior.
// - Control what strings are extracted and how the output is generated.
//
// ### Methods:
// - `DefaultConfig`: Returns a default configuration with reasonable defaults.
// - `WithVerbose`: Enables or disables verbose logging.
// - `WithLogger`: Sets a custom logger instance.
// - `WithConfig`: Applies a predefined configuration.
// - `WithExclude`: Specifies files or directories to exclude from parsing.
// - `WithExtractAll`: Enables or disables extraction of all string literals.
// - `WithHeaderConfig`: Sets the PO file header configuration.
// - `WithHeaderOptions`: Adds header options for PO file generation.
// - `WithHeader`: Sets a custom header for the PO file.
// - `WithFuzzyMatch`: Enables or disables fuzzy matching for translations.
type Config struct {
	Exclude         []string
	ExtractAll      bool
	HeaderConfig    *types.HeaderConfig
	HeaderOptions   []types.HeaderOption
	Header          *types.Header
	Logger          *log.Logger
	Verbose         bool
	CleanDuplicates bool
}

func DefaultConfig(opts ...Option) Config {
	c := Config{
		Header: func() *types.Header {
			h := types.DefaultHeader()
			return &h
		}(),
		Logger:          log.New(io.Discard, "", 0),
		CleanDuplicates: true,
	}

	for _, opt := range opts {
		opt(&c)
	}

	return c
}

type Option func(c *Config)

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

func WithHeaderConfig(h *types.HeaderConfig) Option {
	return func(c *Config) { c.HeaderConfig = h }
}

func WithHeaderOptions(hopts ...types.HeaderOption) Option {
	return func(c *Config) { c.HeaderOptions = hopts }
}

func WithHeader(h *types.Header) Option {
	return func(c *Config) { c.Header = h }
}
