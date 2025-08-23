package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var use = "msgocat"

var root = &cobra.Command{
	Aliases: []string{"msgcat", "cat"},
	Use:     use,
	Args: func(cmd *cobra.Command, args []string) error {
		if filesFrom != "" || directory != "" {
			return nil
		}
		return cobra.MinimumNArgs(1)(cmd, args)
	},
	Short: `Usage: msgocat [OPTION] [INPUTFILE]...

Concatenates and merges the specified PO files.
Find messages which are common to two or more of the specified PO files.`,
	Long: `Usage: msgcat [OPTION] [INPUTFILE]...

Concatenates and merges the specified PO files.
Find messages which are common to two or more of the specified PO files.
By using the --more-than option, greater commonality may be requested
before messages are printed.  Conversely, the --less-than option may be
used to specify less commonality before messages are printed (i.e.
--less-than=2 will only print the unique messages).  Translations,
comments, extracted comments, and file positions will be cumulated, except
that if --use-first is specified, they will be taken from the first PO file
to define them.

Mandatory arguments to long options are mandatory for short options too.`,
	PreRunE: initCfg,
	RunE:    run,
}

func Execute() {
	err := root.Execute()
	if err != nil {
		os.Exit(1)
	}
}
