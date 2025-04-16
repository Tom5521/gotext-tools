package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Tom5521/gotext-tools/v2/internal/slices"
	"github.com/Tom5521/gotext-tools/v2/internal/util"
	goparse "github.com/Tom5521/gotext-tools/v2/pkg/go/parse"
	krfs "github.com/kr/fs"
)

func processInput(inputFiles []string) (*goparse.Parser, error) {
	if filesFrom != "" {
		files, err := readFilesFrom(filesFrom)
		if err != nil {
			return nil, fmt.Errorf("error reading file %s: %w", filesFrom, err)
		}
		inputFiles = files
	}

	if excludeFile != "" {
		var files []string
		files, err := readFilesFrom(excludeFile)
		if err != nil {
			return nil, fmt.Errorf("error reading file %s: %w", excludeFile, err)
		}

		exclude = append(exclude, files...)
	}
	if directory != "" {
		for i, file := range inputFiles {
			inputFiles[i] = filepath.Join(directory, file)
		}
		for i, file := range exclude {
			exclude[i] = filepath.Join(directory, file)
		}
	}
	var err error
	var parser *goparse.Parser

	stdinIndex := slices.Index(inputFiles, "-")
	if stdinIndex != -1 {
		inputFiles = slices.Delete(inputFiles, stdinIndex, stdinIndex+1)
	}

	files, err := readFiles(inputFiles)
	if err != nil {
		return nil, fmt.Errorf("error reading files: %w", err)
	}

	if stdinIndex != -1 {
		files = append(files, os.Stdin)
	}

	// Make the parser.
	parser, err = goparse.NewParserFromFiles(
		files,
		goparse.WithConfig(GoParserCfg),
	)
	if err != nil {
		return nil, fmt.Errorf("error reading files: %w", err)
	}

	return parser, nil
}

func readFiles(paths []string) ([]*os.File, error) {
	seenFiles := make(map[string]bool)

	var files []*os.File
	for _, path := range paths {
		walker := krfs.Walk(path)
		for walker.Step() {
			if util.ShouldSkipFile(walker, exclude, &seenFiles, logger) {
				continue
			}
			file, err := os.Open(walker.Path())
			if err != nil {
				return nil, err
			}

			files = append(files, file)
		}
	}

	return files, nil
}
