package cmd

import (
	"bytes"
	"errors"
	"os"

	"github.com/spf13/cobra"
)

var use = "msgocat"

var root = &cobra.Command{
	Aliases: []string{"msgcat"},
	Use:     use,
	Args:    cobra.MinimumNArgs(1),
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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if lessThan <= 1 {
			return errors.New("impossible selection criteria specified")
		}
		return nil
	},
	RunE: run,
}

func readFilesFrom(path string) ([]string, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, line := range bytes.Split(file, []byte{'\n'}) {
		files = append(files, string(line))
	}

	return files, nil
}

func Execute() {
	err := root.Execute()
	if err != nil {
		os.Exit(1)
	}
}
