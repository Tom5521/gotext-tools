package parser

import (
	"errors"
	"io"
	"slices"
)

type Config struct {
	Files         []string
	InputContent  []byte
	DefaultDomain string
	Output        string
	OutputDir     string

	FallbackOutput io.Writer // is used in case the output is -
	FallbackInput  io.Reader // is used in case the input is -

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
		Files:         []string{"-"},
		Output:        "-",
		DefaultDomain: "default",
		Language:      "en",
		Nplurals:      2,
	}
}

func (c Config) Validate() (errs []error) {
	if len(c.Files) == 0 {
		errs = append(errs, errors.New("there are no input files"))
	}
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
	if len(c.Files) == 1 {
		if c.Files[0] == "-" && c.FallbackInput == nil {
			errs = append(errs, errors.New("input is \"-\", but no fallback has been loaded"))
		}
	} else if len(c.Files) > 1 {
		if slices.Contains(c.Files, "-") {
			errs = append(errs, errors.New("incompatible sources were specified"))
		}
	}

	return
}
