package cmd

import (
	goparse "github.com/Tom5521/xgotext/pkg/go/parse"
	"github.com/Tom5521/xgotext/pkg/po/compiler"
	"github.com/Tom5521/xgotext/pkg/po/types"
)

var (
	ParserCfg   goparse.Config
	CompilerCfg compiler.Config
	HeadersCfg  types.HeaderConfig
)

func initConfig() {
	HeadersCfg = types.HeaderConfig{
		Nplurals:          nplurals,
		ProjectIDVersion:  packageVersion,
		ReportMsgidBugsTo: msgidBugsAddress,
		Language:          lang,
	}
	ParserCfg = goparse.Config{
		Exclude:      exclude,
		ExtractAll:   extractAll,
		HeaderConfig: &HeadersCfg,
		Logger:       logger,
		Verbose:      verbose,
	}
	CompilerCfg = compiler.Config{
		Logger:          logger,
		ForcePo:         forcePo,
		OmitHeader:      omitHeader,
		PackageName:     packageName,
		CopyrightHolder: copyrightHolder,
		ForeignUser:     foreignUser,
		Title:           title,
		NoLocation:      noLocation,
		AddLocation:     compiler.LocationMode(addLocation),
		MsgstrPrefix:    msgstrPrefix,
		MsgstrSuffix:    msgstrSuffix,
		Verbose:         verbose,
	}
}
