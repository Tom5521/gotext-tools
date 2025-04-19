package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Tom5521/gotext-tools/v2/pkg/po/compile"
	"github.com/Tom5521/gotext-tools/v2/pkg/po/parse"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use: os.Args[0],
	Short: `Generate binary message catalog from textual translation descript
ion.`,
	Long: `Usage: msgfmt [OPTION] filename.po ...
Mandatory arguments to long options are mandatory for short optio
ns too.
Similarly for optional arguments.
If input file is -, standard input is read.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if directory != "" {
			return nil
		}
		return cobra.MinimumNArgs(1)(cmd, args)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initCfg()
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if directory != "" {
			output = filepath.Join(directory, output)
			for i, v := range args {
				args[i] = filepath.Join(directory, v)
			}
			if len(args) == 0 {
				entries, err := os.ReadDir(directory)
				if err != nil {
					return err
				}
				for _, de := range entries {
					if !de.IsDir() && filepath.Ext(de.Name()) == ".po" {
						args = append(args, filepath.Join(directory, de.Name()))
					}
				}
			}
		}

		for _, arg := range args {
			path, file := filepath.Split(arg)
			basename := strings.TrimSuffix(file, filepath.Ext(file))
			moPath := filepath.Join(path, basename+".mo")
			poPath := filepath.Join(path, file)

			poFile, err := parse.Po(poPath)
			if err != nil {
				return err
			}

			err = compile.MoToFile(poFile, moPath, compile.MoWithConfig(compilerCfg))
			if err != nil {
				return err
			}

		}

		return
	},
}

func Execute() {
	err := root.Execute()
	if err != nil {
		os.Exit(1)
	}
}
