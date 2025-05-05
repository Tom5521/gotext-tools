package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docs = &cobra.Command{
	Use:   "docs",
	Short: "Generates the documentation of the specified command.",
}

func init() {
	pflags := docs.PersistentFlags()
	pflags.StringP("output", "o", "-", "The output file/directory depending on the command.")
	pflags.StringP(
		"command",
		"c",
		"",
		"Specifies from which command the documentation will be generated.",
	)
	docs.MarkFlagRequired("output")
	docs.MarkFlagRequired("command")

	docs.AddCommand(markdown, man)
}

func genRunners(cmd *cobra.Command) (*os.File, *cobra.Command, error) {
	output, _ := cmd.Flags().GetString("output")
	command, _ := cmd.Flags().GetString("command")

	var toDescribe *cobra.Command
	commands := cmd.Parent().Parent().Commands()
	for _, c := range commands {
		if c.Use == command {
			toDescribe = c
			break
		}
	}

	if toDescribe == nil {
		return nil, nil, fmt.Errorf("the %q command does not exist", command)
	}

	var file *os.File
	var err error
	if output != "-" {
		file, err = os.Create(output)
		if err != nil {
			return nil, nil, err
		}
	} else {
		file = os.Stdout
	}

	return file, toDescribe, err
}

var (
	markdown = &cobra.Command{
		Use:   "markdown",
		Short: "Generate documentation in markdown.",
		RunE: func(cmd *cobra.Command, args []string) error {
			file, d, err := genRunners(cmd)
			if err != nil {
				return err
			}
			if file != os.Stdout {
				defer file.Close()
			}

			return doc.GenMarkdown(d, file)
		},
	}
	man = &cobra.Command{
		Use:   "man",
		Short: "Generate documentation for man.",
		RunE: func(cmd *cobra.Command, args []string) error {
			file, d, err := genRunners(cmd)
			if err != nil {
				return err
			}

			if file != os.Stdout {
				defer file.Close()
			}

			return doc.GenMan(d, nil, file)
		},
	}
)
