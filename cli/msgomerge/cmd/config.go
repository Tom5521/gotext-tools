package cmd

import (
	"os"

	"github.com/Tom5521/gotext-tools/v2/pkg/po"
	"github.com/Tom5521/gotext-tools/v2/pkg/po/compile"
	"golang.org/x/term"
)

var (
	mergeCfg    po.MergeConfig
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

	if term.IsTerminal(int(os.Stdout.Fd())) {
		compilerCfg.Highlight = compile.DefaultHighligh
	}

	mergeCfg = po.MergeConfig{
		FuzzyMatch: !noFuzzyMatching,
		Sort:       true,
	}
}
