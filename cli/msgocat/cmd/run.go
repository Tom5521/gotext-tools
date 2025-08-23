package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/Tom5521/gotext-tools/v2/pkg/po"
	"github.com/Tom5521/gotext-tools/v2/pkg/po/compile"
	"github.com/Tom5521/gotext-tools/v2/pkg/po/parse"
	"github.com/spf13/cobra"
)

// run is the main execution function for the command.
func run(cmd *cobra.Command, args []string) error {
	args, err := expandArgsWithFiles(args)
	if err != nil {
		return err
	}

	err = cobra.MinimumNArgs(1)(cmd, args)
	if err != nil {
		return err
	}

	files := make([]*po.File, len(args))
	for i, filename := range args {
		var file *po.File
		if filename == "-" {
			file, err = parse.PoFromReader(os.Stdin, "-")
		} else {
			file, err = parse.Po(filename)
		}
		if err != nil {
			return fmt.Errorf("error parsing PO file (%s): %w", filename, err)
		}
		errs := file.Validate()
		if len(errs) > 0 {
			return fmt.Errorf("invalid PO file (%s): %w", filename, errs[0])
		}

		files[i] = file
	}

	mergedFiles := po.MsgcatMergeFiles(
		files,
		po.MsgcatMergeWithConfig(mergeCfg),
	)

	if lang != "" {
		hindex := mergedFiles.Index("", "")
		if hindex != -1 {
			oldHeader := mergedFiles[hindex]
			header := mergedFiles.HeaderFromIndex(hindex)
			header.Set("Language", lang)

			newHeader := header.ToEntry()
			newHeader.Flags = oldHeader.Flags
			newHeader.Comments = oldHeader.Comments
			newHeader.ExtractedComments = oldHeader.ExtractedComments
			newHeader.Previous = oldHeader.Previous

			mergedFiles[hindex] = newHeader
		}
	}

	if sortOutput {
		mergedFiles = mergedFiles.Sort()
	}
	if sortByFile {
		mergedFiles = mergedFiles.SortFunc(po.CompareEntryByFile)
	}
	var out io.Writer

	if output == "-" {
		out = os.Stdout
	} else {
		var fileRef *os.File
		fileRef, err = os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
		if err != nil {
			return fmt.Errorf("error openning output file %s: %w", output, err)
		}
		defer fileRef.Close()
		out = fileRef
	}

	return compile.PoToWriter(
		mergedFiles,
		out,
		compile.PoWithConfig(compilerCfg),
	)
}
