package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	goparse "github.com/Tom5521/xgotext/pkg/go/parse"
	"github.com/Tom5521/xgotext/pkg/po/compiler"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "xgotext",
	Short: "Extract translatable strings from given input files.",
	Long: `Extract translatable strings from given input files.

Mandatory arguments to long options are mandatory for short options too.
Similarly for optional arguments.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
	},
	RunE: func(cmd *cobra.Command, inputfiles []string) (err error) {
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

		p, err := goparse.NewParserFromFiles(
			inputfiles,
			ParserCfg,
		)
		p.HeaderCfg = &HeadersCfg
		if err != nil {
			return fmt.Errorf("error parsing files: %w", err)
		}

		pofile := p.Parse()
		if len(p.Errors()) > 0 {
			return fmt.Errorf("errors in entries parsing (%d): %w", len(p.Errors()), p.Errors()[0])
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

			flags := os.O_RDWR

			if os.IsExist(err) && !forcePo && output != "" && !joinExisting {
				return fmt.Errorf("file %s already exists", outputFile)
			} else if os.IsNotExist(err) {
				flags |= os.O_CREATE
			} else if err != nil {
				return fmt.Errorf("error getting file %s information: %w", outputFile, err)
			}

			file, err = os.OpenFile(outputFile, flags, os.ModePerm)
			if err != nil {
				return fmt.Errorf("error opening file %s: %w", outputFile, err)
			}
			defer file.Close()

			if joinExisting {
				return join(p, file)
			}

			if stat.Size() != 0 {
				err = file.Truncate(0)
				if err != nil {
					return fmt.Errorf("error truncating file %s: %w", outputFile, err)
				}
			}
			out = file
		}

		compiler := compiler.Compiler{
			File:   pofile,
			Config: CompilerCfg,
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
