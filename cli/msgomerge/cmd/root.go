package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   os.Args[0],
	Short: "Merges two Uniforum style .po files together.",
	Long: `Merges two Uniforum style .po files together.  The def.po file is an
existing PO file with translations which will be taken over to the newly
created file as long as they still match; comments will be preserved,
but extracted comments and file positions will be discarded.  The ref.pot
file is the last created PO file with up-to-date source references but
old translations, or a PO Template file (generally created by xgettext);
any translations or comments in the file will be discarded, however dot
comments and file positions will be preserved.  Where an exact match
cannot be found, fuzzy matching is used to produce better results.`,
	Args:   cobra.RangeArgs(2, 2),
	PreRun: func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	err := root.Execute()
	if err != nil {
		os.Exit(1)
	}
}
