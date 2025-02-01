package cmd

import (
	"log"
	"os"

	"github.com/Tom5521/xgotext/pkg/po/config"
)

var cfg = config.Config{
	Logger:           log.New(os.Stdout, "LOG: ", log.Ltime),
	DefaultDomain:    defaultDomain,
	ForcePo:          forcePo,
	NoLocation:       noLocation,
	AddLocation:      addLocation,
	OmitHeader:       omitHeader,
	PackageName:      packageName,
	PackageVersion:   packageVersion,
	Language:         lang,
	Nplurals:         nplurals,
	Exclude:          exclude,
	ForeignUser:      foreignUser,
	MsgidBugsAddress: msgidBugsAddress,
	Title:            title,
	CopyrightHolder:  copyrightHolder,
	JoinExisting:     joinExisting,
	ExtractAll:       extractAll,
	Verbose:          verbose,
}

func init() {
	cfg.Msgstr.Prefix = msgstrPrefix
	cfg.Msgstr.Suffix = msgstrSuffix
}
