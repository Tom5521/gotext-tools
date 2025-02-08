package compiler

import "errors"

type Config struct {
	ForcePo         bool
	OmitHeader      bool
	PackageName     string
	CopyrightHolder string
	ForeignUser     bool
	Title           string
	NoLocation      bool
	AddLocation     string
	MsgstrPrefix    string
	MsgstrSuffix    string
}

func (c Config) Validate() error {
	if c.NoLocation && c.AddLocation != "never" {
		return errors.New("noLocation and AddLocation are in conflict")
	}

	return nil
}
