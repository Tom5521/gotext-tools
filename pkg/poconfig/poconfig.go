package poconfig

import (
	"errors"
)

type Config struct {
	DefaultDomain  string
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
		DefaultDomain: "default",
		Language:      "en",
		Nplurals:      2,
	}
}

func (c Config) Validate() (errs []error) {
	if c.NoLocation && c.AddLocation != "never" {
		errs = append(errs, errors.New("noLocation and AddLocation are in conflict"))
	}

	if c.Nplurals == 0 {
		errs = append(errs, errors.New("nplurals is equal to 0"))
	}

	return
}
