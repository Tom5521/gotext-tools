package cmd

import (
	"github.com/Tom5521/gotext-tools/pkg/po"
	"github.com/Tom5521/gotext-tools/pkg/po/compiler"
)

var (
	mergeCfg    po.MergeConfig
	headerCfg   po.HeaderConfig
	compilerCfg compiler.PoConfig
)

func initConfig() {
	compilerCfg = compiler.PoConfig{
		NoLocation:  noLocation,
		AddLocation: compiler.PoLocationMode(addLocation),
		WordWrap:    !noWrap,
		ForcePo:     forcePo,
		OmitHeader:  true,
	}
	mergeCfg = po.MergeConfig{
		FuzzyMatch: !noFuzzyMatching,
		Sort:       true,
	}
}
