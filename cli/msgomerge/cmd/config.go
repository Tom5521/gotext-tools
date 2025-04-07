package cmd

import (
	"github.com/Tom5521/xgotext/pkg/po"
	"github.com/Tom5521/xgotext/pkg/po/compiler"
)

var (
	mergeCfg    po.MergeConfig
	headerCfg   po.HeaderConfig
	compilerCfg compiler.PoConfig
)

func initConfig() {
	headerCfg = po.HeaderConfig{
		Language: lang,
	}
	compilerCfg = compiler.PoConfig{
		NoLocation:  noLocation,
		AddLocation: compiler.PoLocationMode(addLocation),
	}
	mergeCfg = po.MergeConfig{
		FuzzyMatch: !noFuzzyMatching,
	}
}
