package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/Tom5521/xgotext/pkg/goparse"
	"github.com/Tom5521/xgotext/pkg/po/compiler"
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
		if filesFrom != "" {
			var files []string
			files, err = readFilesFrom(filesFrom)
			if err != nil {
				return fmt.Errorf("error reading file %s: %w", filesFrom, err)
			}
			inputfiles = append(inputfiles, files...)
		}
		if excludeFile != "" {
			var files []string
			files, err = readFilesFrom(excludeFile)
			if err != nil {
				return fmt.Errorf("error reading file %s: %w", excludeFile, err)
			}
			exclude = append(exclude, files...)
		}

		if directory != "" {
			for i, file := range inputfiles {
				inputfiles[i] = filepath.Join(directory, file)
			}
			for i, file := range exclude {
				exclude[i] = filepath.Join(directory, file)
			}
		}

		config := config.Config{
			Logger:           log.New(os.Stdout, "INFO: ", log.Ltime),
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
			Verbose:          verbose,
		}
		config.Msgstr.Prefix = msgstrPrefix
		config.Msgstr.Suffix = msgstrSuffix

		p, err := goparse.NewParserFromFiles(
			inputfiles,
			config,
		)
		if err != nil {
			return fmt.Errorf("error parsing files: %w", err)
		}

		translations, errs := p.Parse()
		if len(errs) > 0 {
			return fmt.Errorf("errors in translations parsing (%d): %w", len(errs), errs[0])
		}

		outputFile := filepath.Join(outputDir, defaultDomain+".pot")

		var out io.Writer

		switch {
		case output == "-":
			out = os.Stdout
		case output != "":
			if outputDir != "" {
				output = filepath.Join(outputDir, output)
			}
			outputFile = output
			fallthrough
		default:
			var file *os.File
			var stat os.FileInfo
			stat, err = os.Stat(outputFile)
			if err != nil {
				return fmt.Errorf("error getting file %s information: %w", outputFile, err)
			}
			if os.IsExist(err) && !forcePo && output != "" {
				return fmt.Errorf("file %s already exists", outputFile)
			}

			file, err = os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE, os.ModePerm)
			if err != nil {
				return fmt.Errorf("error opening file %s: %w", outputFile, err)
			}
			defer file.Close()

			if stat.Size() != 0 {
				err = file.Truncate(0)
				if err != nil {
					return fmt.Errorf("error truncating file %s: %w", outputFile, err)
				}
			}
			out = file
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
