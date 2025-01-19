package poconfig

import (
	"errors"
	"io"
)

type Config struct {
	DefaultDomain string
	Output        string
	OutputDir     string

	FallbackOutput io.Writer // is used in case the output is -

	ForcePo        bool
	NoLocation     bool
	AddLocation    string
	OmitHeader     bool
	PackageName    string
	PackageVersion string
	Language       string
	Nplurals       uint
	Exclude        []string
	JoinExisting   bool
	ExtractAll     bool
	Msgstr         struct {
		Prefix string
		Suffix string
	}
}

func DefaultConfig() Config {
	return Config{
		Output:        "-",
		DefaultDomain: "default",
		Language:      "en",
		Nplurals:      2,
	}
}

func (c Config) Validate() (errs []error) {
	if c.Output == "" {
		errs = append(errs, errors.New("there must be an output"))
	}

	if c.Output != "" && c.OutputDir != "" {
		errs = append(errs, errors.New("output and OutputDir options are mutually exclusive"))
	}

	if c.Output == "" && c.OutputDir == "" {
		errs = append(errs, errors.New("there are no outputs"))
	}

	if c.NoLocation && c.AddLocation != "never" {
		errs = append(errs, errors.New("noLocation and AddLocation are in conflict"))
	}

	if c.Nplurals == 0 {
		errs = append(errs, errors.New("nplurals is equal to 0"))
	}

	if c.Output == "-" && c.FallbackOutput == nil {
		errs = append(errs, errors.New("output is \"-\", but no fallback has been loaded"))
	}

	return
}
