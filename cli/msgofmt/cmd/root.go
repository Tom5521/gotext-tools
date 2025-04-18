package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Tom5521/gotext-tools/v2/pkg/po"
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
				if v == "-" {
					continue
				}
				args[i] = filepath.Join(directory, v)
			}
			if len(args) >= 0 {
				var entries []os.DirEntry
				entries, err = os.ReadDir(directory)
				if err != nil {
					return
				}
				if len(args) > 0 {
					for _, de := range entries {
						if !de.IsDir() && filepath.Ext(de.Name()) == ".po" {
							args = append(args, filepath.Join(directory, de.Name()))
						}
					}
				} else {
					for _, de := range entries {
						if de.IsDir() {
							continue
						}
						basename := strings.TrimSuffix(de.Name(), filepath.Ext(de.Name()))
						newMo := filepath.Join(directory, basename+".mo")

						var poFile *po.File
						poFile, err = parse.Po(filepath.Join(directory, de.Name()))
						if err != nil {
							return
						}
						if err = poFile.Validate(); err != nil {
							return
						}
						err = compile.MoToFile(poFile, newMo)
						if err != nil {
							return
						}
					}
					return
				}
			}
		}

		var allEntries po.Entries
		for _, arg := range args {
			var poFile *po.File
			if arg == "-" {
				poFile, err = parse.PoFromReader(os.Stdin, "stdin")
			} else {
				poFile, err = parse.Po(arg)
			}
			if err != nil {
				return
			}

			allEntries = append(allEntries, poFile.Entries...)
		}

		if err = allEntries.Validate(); err != nil {
			return
		}

		return compile.MoToFile(allEntries, output, compile.MoWithConfig(compilerCfg))
	},
}

func Execute() {
	err := root.Execute()
	if err != nil {
		os.Exit(1)
	}
}
