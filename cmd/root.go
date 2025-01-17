package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "xgotext",
	Short: "Extract translatable strings from given input files.",
	Long: `Extract translatable strings from given input files.

Mandatory arguments to long options are mandatory for short options too.
Similarly for optional arguments.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func Execute() {
	err := root.Execute()
	if err != nil {
		os.Exit(1)
	}
}
