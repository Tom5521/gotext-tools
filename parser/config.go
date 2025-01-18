package parser

import (
	"fmt"
)

type Config struct {
	Files          []string
	DefaultDomain  string
	Output         string
	OutputDir      string
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
		errs = append(errs, fmt.Errorf("there are no input files"))
	}
	if c.Output == "" {
		errs = append(errs, fmt.Errorf("there must be an output"))
	}

	if c.Output != "" && c.OutputDir != "" {
		errs = append(errs, fmt.Errorf("Output and OutputDir options are mutually exclusive"))
	}

	if c.Output == "" && c.OutputDir == "" {
		errs = append(errs, fmt.Errorf("there are no outputs"))
	}

	if c.NoLocation && c.AddLocation != "never" {
		errs = append(errs, fmt.Errorf("NoLocation and AddLocation are in conflict"))
	}

	if c.Nplurals == 0 {
		errs = append(errs, fmt.Errorf("Npluras is equal to 0"))
	}

	return
}
