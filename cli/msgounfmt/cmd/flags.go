package cmd

import (
	"os"

	"github.com/Tom5521/gotext-tools/v2/pkg/po/compile"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	output       string
	color        string
	forcePo      bool
	noWrap       bool
	sortOutput   bool
	ignoreErrors bool
)

func init() {
	flag := root.Flags()
	flag.StringVarP(&output, "output", "o", "-", `write
output to specified file
The results are written to standard output if no output file is specified
or if it is -.`)
	flag.StringVar(&color, "color", "auto", `use colors and other text attributes if WHEN.
WHEN may be 'always', 'never' or 'auto'.`)

	flag.BoolVarP(&sortOutput, "sort-output", "s", false, `generate sorted output`)
	flag.BoolVar(&noWrap, "no-wrap", false, `do not break long message lines, longer than
the output page width, into several lines`)
	flag.BoolVarP(&forcePo, "force-po", "f", false, `write PO file even if empty`)
	flag.BoolVar(
		&ignoreErrors,
		"ignore-errors",
		false,
		"skip non-critical errors in the data such as duplicate and/or unsorted entries.",
	)
}

var compilerCfg = compile.DefaultPoConfig()

func initCfg(cmd *cobra.Command, args []string) error {
	compilerCfg.WordWrap = !noWrap
	compilerCfg.ForcePo = forcePo

	switch color {
	case "auto":
		if !term.IsTerminal(int(os.Stdout.Fd())) || output != "-" {
			break
		}
		fallthrough
	case "always":
		compilerCfg.Highlight = compile.DefaultHighlight
	}

	return nil
}
