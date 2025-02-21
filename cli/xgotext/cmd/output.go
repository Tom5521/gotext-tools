package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

func processOutput() (*os.File, error) {
	if output == "-" {
		return os.Stdout, nil
	}

	customOutput := output != ""
	outputFilePath := filepath.Join(outputDir, defaultDomain+".pot")
	if customOutput {
		outputFilePath = filepath.Join(outputDir, output)
	}

	_, err := os.OpenFile(outputFilePath, os.O_RDWR, os.ModePerm)
	if os.IsExist(err) && !forcePo && customOutput {
		return nil, fmt.Errorf("file %s already exists", outputFilePath)
	}

	// Truncate file.
	return os.Create(outputFilePath)
}
