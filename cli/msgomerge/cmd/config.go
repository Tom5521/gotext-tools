package cmd

import (
	"log"
	"os"

	"github.com/Tom5521/gotext-tools/v2/pkg/po"
	"github.com/Tom5521/gotext-tools/v2/pkg/po/compile"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	mergeCfg    po.MergeConfig
	compilerCfg compile.PoConfig
)

func initConfig(cmd *cobra.Command, args []string) {
	compilerCfg = compile.PoConfig{
		NoLocation:  noLocation,
		AddLocation: compile.PoLocationMode(addLocation),
		WordWrap:    !noWrap,
		ForcePo:     forcePo,
		OmitHeader:  true,
		Verbose:     verbose,
		Logger:      log.Default(),
	}

	switch color {
	case "auto":
		if !term.IsTerminal(int(os.Stdout.Fd())) || outputPath != "-" {
			break
		}
		fallthrough
	case "always":
		compilerCfg.Highlight = compile.DefaultHighlight
	}

	mergeCfg = po.MergeConfig{
		FuzzyMatch: !noFuzzyMatching,
		Sort:       true,
	}
}
