package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/Tom5521/gotext-tools/v2/pkg/po"
	"github.com/Tom5521/gotext-tools/v2/pkg/po/parse"
)

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

// expandArgsWithFiles handles the --files-from functionality.
func expandArgsWithFiles(args []string) ([]string, error) {
	if filesFrom == "" {
		return args, nil
	}

	files, err := readFilesFrom(filesFrom)
	if err != nil {
		return nil, err
	}
	return append(args, files...), nil
}

// parsePOFile parses a PO file from disk.
func parsePOFile(filename string) (*po.File, error) {
	fileRef, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error openning file: %w", err)
	}
	defer fileRef.Close()

	poFile, err := parse.PoFromFile(fileRef, parse.PoWithCleanDuplicates(false))
	if err != nil {
		return nil, fmt.Errorf("error parsing file: %w", err)
	}

	if errs := poFile.Validate(); len(errs) > 0 {
		return nil, errs[0]
	}
	return poFile, nil
}
