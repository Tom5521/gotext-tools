package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/Tom5521/gotext-tools/v2/pkg/po/compile"
	"github.com/spf13/cobra"
)

var logger = log.New(os.Stdout, "", log.Ltime)

var root = &cobra.Command{
	Use:   os.Args[0],
	Short: "Extract translatable strings from given input files.",
	Long: `Extract translatable strings from given input files.
Mandatory arguments to long options are mandatory for short options too.
Similarly for optional arguments.`,
	Example: fmt.Sprintf(`%s -a my-go-proyect/ -o -
%s -a my-go-proyect/ -o en.pot -lang en
%s main.go -o main.pot -lang en`,
		os.Args[0], os.Args[0], os.Args[0],
	),
	PreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
	},
	RunE: func(cmd *cobra.Command, inputfiles []string) (err error) {
		parser, err := processInput(inputfiles)
		if err != nil {
			return
		}

		parsedFile := parser.Parse()
		if len(parser.Errors()) > 0 {
			return fmt.Errorf(
				"errors in entries parsing (%d): %w",
				len(parser.Errors()),
				parser.Errors()[0],
			)
		}

		out, err := processOutput()
		if err != nil {
			return err
		}
		defer func() {
			if out != os.Stdout {
				out.Close()
			}
		}()

		if joinExisting {
			return join(parser, out)
		}

		err = compile.PoToWriter(parsedFile, out, compile.PoWithConfig(CompilerCfg))
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
