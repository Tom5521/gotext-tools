package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docTree = &cobra.Command{
	Use:   "doc-tree",
	Short: "Generate markdown tree.",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			return err
		}

		if _, err = os.Stat(dir); os.IsNotExist(err) {
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return err
			}
		}

		return doc.GenMarkdownTree(cmd.Parent(), dir)
	},
}

func init() {
	docTree.Flags().
		StringP("dir", "D", "", "specifies in which directory the documentation will be written")
}
