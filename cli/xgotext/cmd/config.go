package cmd

import (
	goparse "github.com/Tom5521/xgotext/pkg/go/parse"
	"github.com/Tom5521/xgotext/pkg/po"
	"github.com/Tom5521/xgotext/pkg/po/compiler"
	poparse "github.com/Tom5521/xgotext/pkg/po/parse"
)

var (
	PoParserCfg poparse.PoConfig
	GoParserCfg goparse.Config
	CompilerCfg compiler.PoConfig
	HeadersCfg  po.HeaderConfig
)

func initConfig() {
	HeadersCfg = po.DefaultHeaderConfig()
	HeadersCfg.Nplurals = nplurals
	HeadersCfg.ProjectIDVersion = packageVersion
	HeadersCfg.ReportMsgidBugsTo = msgidBugsAddress
	HeadersCfg.Language = lang

	GoParserCfg = goparse.Config{
		Exclude:      exclude,
		ExtractAll:   extractAll,
		HeaderConfig: &HeadersCfg,
		Logger:       logger,
		Verbose:      verbose,
	}
	CompilerCfg = compiler.PoConfig{
		Logger:          logger,
		ForcePo:         forcePo,
		OmitHeader:      omitHeader,
		PackageName:     packageName,
		CopyrightHolder: copyrightHolder,
		ForeignUser:     foreignUser,
		Title:           title,
		NoLocation:      noLocation,
		AddLocation:     compiler.PoLocationMode(addLocation),
		MsgstrPrefix:    msgstrPrefix,
		MsgstrSuffix:    msgstrSuffix,
		Verbose:         verbose,
		HeaderComments:  true,
		HeaderFields:    true,
	}
	PoParserCfg = poparse.PoConfig{
		Logger: logger,
	}
}
