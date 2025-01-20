package cmd

import (
	"bytes"
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
			return
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
				return
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
			return
		}

		translations, errs := p.Parse()
		if len(errs) > 0 {
			return errs[0]
		}

		err = os.Chdir(startDir)
		if err != nil {
			return
		}

		outputFile := filepath.Join(outputDir, defaultDomain+".pot")

		var out io.Writer

		if output != "" {
			if output == "-" {
				out = os.Stdout
			} else {
				out, err = os.OpenFile(output, os.O_WRONLY|os.O_CREATE, os.ModePerm)
			}
		} else {
			out, err = os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE, os.ModePerm)
		}
		if err != nil {
			return
		}

		compiler := compiler.Compiler{
			Translations: translations,
			Config:       config,
		}

		err = compiler.CompileToWriter(out)

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
