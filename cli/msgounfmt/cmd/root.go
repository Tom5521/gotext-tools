package cmd

import (
	"os"

	"github.com/Tom5521/gotext-tools/v2/pkg/po"
	"github.com/Tom5521/gotext-tools/v2/pkg/po/compile"
	"github.com/Tom5521/gotext-tools/v2/pkg/po/parse"
	"github.com/spf13/cobra"
)

const use = "msgounfmt"

var root = &cobra.Command{
	Use:   use,
	Short: `Convert binary message catalog to Uniforum style .po file.`,
	Long: `Usage: msgounfmt [OPTIONS]... [FILE]...

Convert binary message catalog to Uniforum style .po file.

Mandatory arguments to long options are mandatory for short options too.

Input file location:
  FILE ...                    input .mo files
If no input file is given or if it is -, standard input is read.`,
	PreRunE: initCfg,
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := &po.File{
			Name: "input.mo",
		}
		for _, arg := range args {
			var f *po.File
			var err error
			if arg == "-" {
				f, err = parse.MoFromReader(os.Stdin, "stdin")
			} else {
				f, err = parse.Mo(arg)
			}

			if err != nil {
				return err
			}

			err = f.Validate()
			if err != nil && !ignoreErrors {
				return err
			}

			file.Entries = append(file.Entries, f.Entries...)
		}

		if sortOutput {
			file.Entries = file.SortFunc(po.CompareEntryByID)
		}

		if output == "-" {
			return compile.PoToWriter(file, os.Stdout, compile.PoWithConfig(compilerCfg))
		} else {
			return compile.PoToFile(file, output, compile.PoWithConfig(compilerCfg))
		}
	},
}

func Execute() {
	err := root.Execute()
	if err != nil {
		os.Exit(1)
	}
}
