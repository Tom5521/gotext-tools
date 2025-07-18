package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Tom5521/gotext-tools/v2/pkg/po"
	"github.com/Tom5521/gotext-tools/v2/pkg/po/compile"
	"github.com/Tom5521/gotext-tools/v2/pkg/po/parse"

	"github.com/spf13/cobra"
)

var use = "msgomerge"

var root = &cobra.Command{
	Use:   use,
	Short: "Merges two Uniforum style .po files together.",
	Long: `Merges two Uniforum style .po files together.  
The def.po file is an existing PO file with translations which will be
taken over to the newly created file as long as they still match;
comments will be preserved, but extracted comments and file positions will
be discarded.
The ref.pot file is the last created PO file with up-to-date source references but
old translations, or a PO Template file (generally created by xgettext);
any translations or comments in the file will be discarded, however dot
comments and file positions will be preserved.
Where an exact match cannot be found, fuzzy matching is used to produce better results.`,
	Example: fmt.Sprintf(`%s es.po en.pot -o -
%s es.po en.pot -o merged.po
%s es.po en.pot -U`,
		use, use, use,
	),
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 && len(compendium) < 1 {
			return cobra.ExactArgs(2)(cmd, args)
		}

		return nil
	},
	PreRun: initConfig,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		defPath, refPath := args[0], args[1]

		if directory != "" {
			for i, v := range compendium {
				compendium[i] = filepath.Join(directory, v)
			}
			defPath = filepath.Join(directory, defPath)
			refPath = filepath.Join(directory, refPath)
		}

		var defFile, refFile *os.File
		if defFile, err = os.OpenFile(defPath, os.O_RDWR, os.ModePerm); err != nil {
			return err
		}
		if refFile, err = os.OpenFile(refPath, os.O_RDWR, os.ModePerm); err != nil {
			return err
		}
		defer defFile.Close()
		defer refFile.Close()

		var outWriter io.Writer
		if outputPath == "-" {
			outWriter = os.Stdout
		} else {
			var file *os.File
			file, err = os.OpenFile(outputPath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
			if err != nil {
				return err
			}
			defer file.Close()
			outWriter = file
		}

		var def, ref *po.File
		if def, err = parse.PoFromFile(defFile); err != nil {
			return err
		}
		if ref, err = parse.PoFromFile(refFile); err != nil {
			return err
		}

		// Read compendiums.
		for _, comp := range compendium {
			var c *po.File
			c, err = parse.Po(comp)
			if err != nil {
				return err
			}
			ref.Entries = append(ref.Entries, c.Entries...)
		}
		if len(compendium) > 0 {
			ref.Entries = ref.Solve()
		}

		if update {
			// Reset and truncate defFile
			if _, err := defFile.Seek(0, 0); err != nil {
				return err
			}
			if err := defFile.Truncate(0); err != nil {
				return err
			}

			outWriter = io.MultiWriter(defFile, outWriter)
		}

		out := &po.File{
			Name:    outputPath,
			Entries: po.MergeWithConfig(mergeCfg, def.Entries, ref.Entries),
		}

		return compile.PoToWriter(out, outWriter, compile.PoWithConfig(compilerCfg))
	},
}

func Execute() {
	err := root.Execute()
	if err != nil {
		os.Exit(1)
	}
}
