package cmd

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

func readFilesFrom(content string) []string {
	var files []string
	for _, line := range strings.Split(content, "\n") {
		files = append(files, line)
	}

	return files
}

// expandArgsWithFiles handles the --files-from functionality.
func expandArgsWithFiles(args []string) ([]string, error) {
	var filesFromBinContent []byte
	var filesFromContent []string
	var err error

	if filesFrom == "-" {
		filesFromBinContent, err = io.ReadAll(os.Stdin)
	} else if filesFrom != "" {
		filesFromBinContent, err = os.ReadFile(filesFrom)
	}

	if err != nil {
		return nil, err
	}

	if len(filesFromContent) > 0 {
		filesFromContent = readFilesFrom(string(filesFromBinContent))
	}

	var directoryContent []string

	if directory != "" {
		var entries []os.DirEntry
		entries, err = os.ReadDir(directory)
		if err != nil {
			return nil, err
		}

		for _, de := range entries {
			directoryContent = append(directoryContent, filepath.Join(directory, de.Name()))
		}
	}

	return append(args, append(filesFromContent, directoryContent...)...), nil
}
