package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Tom5521/xgotext/pkg/compiler"
	"github.com/Tom5521/xgotext/pkg/goparse"
	"github.com/Tom5521/xgotext/pkg/po/config"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "xgotext",
	Short: "Extract translatable strings from given input files.",
	Long: `Extract translatable strings from given input files.

Mandatory arguments to long options are mandatory for short options too.
Similarly for optional arguments.`,
	RunE: func(_ *cobra.Command, inputfiles []string) (err error) {
		startDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting current directory: %w", err)
		}

		if directory != "" {
			err = os.Chdir(directory)
			if err != nil {
				return
			}
		}

		if filesFrom != "" {
			var files []string
			files, err = readFilesFrom(filesFrom)
			if err != nil {
				return fmt.Errorf("error reading file %s: %w", filesFrom, err)
			}
			inputfiles = append(inputfiles, files...)
		}

		config := config.Config{
			DefaultDomain:    defaultDomain,
			ForcePo:          forcePo,
			NoLocation:       noLocation,
			AddLocation:      addLocation,
			OmitHeader:       omitHeader,
			PackageName:      packageName,
			PackageVersion:   packageVersion,
			Language:         lang,
			Nplurals:         nplurals,
			Exclude:          exclude,
			ForeignUser:      foreignUser,
			MsgidBugsAddress: msgidBugsAddress,
			Title:            title,
			CopyrightHolder:  copyrightHolder,
			JoinExisting:     joinExisting,
			ExtractAll:       extractAll,
		}
		config.Msgstr.Prefix = msgstrPrefix
		config.Msgstr.Suffix = msgstrSuffix

		p, err := goparse.NewParserFromFiles(inputfiles, config)
		if err != nil {
			return fmt.Errorf("error parsing files: %w", err)
		}

		err = os.Chdir(startDir)
		if err != nil {
			return
		}

		translations, errs := p.Parse()
		if len(errs) > 0 {
			return fmt.Errorf("errors in translations parsing (%d): %w", len(errs), errs[0])
		}

		outputFile := filepath.Join(outputDir, defaultDomain+".pot")

		var out io.Writer

		if output != "" {
			if output == "-" {
				out = os.Stdout
			} else {
				outputFile = output
				out, err = os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE, os.ModePerm)
			}
		} else {
			out, err = os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE, os.ModePerm)
		}
		if err != nil {
			return fmt.Errorf("error oppening file %s: %w", outputFile, err)
		}
		if w, ok := out.(interface{ Truncate(int64) error }); ok {
			err = w.Truncate(0)
			if err != nil {
				return fmt.Errorf("error truncating file %s: %w", outputFile, err)
			}
		}

		compiler := compiler.Compiler{
			Translations: translations,
			Config:       config,
		}

		err = compiler.CompileToWriter(out)
		if err != nil {
			return fmt.Errorf("error compiling translations: %w", err)
		}

		return
	},
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
