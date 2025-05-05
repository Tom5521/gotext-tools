package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "gotext-tools",
	Short: "A wrapper for the CLI tools from github.com/Tom5521/gotext-tools/v2/cli",
	Long:  "A wrapper for the CLI tools from github.com/Tom5521/gotext-tools/v2/cli, its main function is to avoid occupying commands in the PATH.",
}

func init() {
	root.AddCommand(msgofmt, msgomerge, xgotext, docs)
}

func Execute() {
	err := root.Execute()
	if err != nil {
		os.Exit(1)
	}
}
