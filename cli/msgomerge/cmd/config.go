package cmd

import (
	"github.com/Tom5521/gotext-tools/pkg/po"
	"github.com/Tom5521/gotext-tools/pkg/po/compile"
)

var (
	mergeCfg    po.MergeConfig
	headerCfg   po.HeaderConfig
	compilerCfg compile.PoConfig
)

func initConfig() {
	compilerCfg = compile.PoConfig{
		NoLocation:  noLocation,
		AddLocation: compile.PoLocationMode(addLocation),
		WordWrap:    !noWrap,
		ForcePo:     forcePo,
		OmitHeader:  true,
	}
	mergeCfg = po.MergeConfig{
		FuzzyMatch: !noFuzzyMatching,
		Sort:       true,
	}
}
