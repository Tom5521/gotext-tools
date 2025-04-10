package cmd

import (
	"os"
	"path/filepath"

	"github.com/Tom5521/xgotext/pkg/po"
	"github.com/Tom5521/xgotext/pkg/po/compiler"
	"github.com/Tom5521/xgotext/pkg/po/parse"
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
	Args: cobra.ExactArgs(2),
	PreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		defPath, refPath := args[0], args[1]

		if directory != "" {
			for i, v := range compendium {
				compendium[i] = filepath.Join(directory, v)
			}
			defPath = filepath.Join(directory, defPath)
			refPath = filepath.Join(directory, refPath)
		}

		var outFile *os.File
		if outputPath == "-" {
			outFile = os.Stdout
		} else {
			outFile, err = os.OpenFile(outputPath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
			if err != nil {
				return err
			}
			defer outFile.Close()
		}
		var def, ref *po.File
		// Read files.
		{
			def, err = parse.ParsePo(defPath)
			if err != nil {
				return err
			}
			ref, err = parse.ParsePo(refPath)
			if err != nil {
				return err
			}
		}
		// Read compendiums.
		{
			for _, comp := range compendium {
				c, err := parse.ParsePo(comp)
				if err != nil {
					return err
				}
				ref.Entries = append(ref.Entries, c.Entries...)
			}
			if len(compendium) > 0 {
				ref.Entries = ref.Solve()
			}
		}

		out := &po.File{
			Name:    outputPath,
			Entries: po.MergeWithConfig(mergeCfg, def.Entries, ref.Entries),
		}

		comp := compiler.PoCompiler{
			File:   out,
			Config: compilerCfg,
		}
		return comp.ToWriter(outFile)
	},
}

func Execute() {
	err := root.Execute()
	if err != nil {
		os.Exit(1)
	}
}
