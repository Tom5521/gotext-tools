package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"

	"github.com/gookit/color"
	"github.com/spf13/pflag"
)

const (
	baseRegex    = `%s\.Get\s*\(\s*.*\s*\)`
	importsRegex = `import\s+(\([^\)]*\)|[^\(\)\s]+)`
	wantedPkg    = `"github.com/leonelquinteros/gotext"`
)

var importsCompiler = regexp.MustCompile(importsRegex)

func main() {
	pflag.Parse()

	msgids, err := processPaths(input, exclude)
	if err != nil {
		color.Errorln(err)
		os.Exit(1)
	}
	err = writeOutput(ExportMsgIDs(msgids))
	if err != nil {
		color.Errorln(err)
		os.Exit(1)
	}
}

func writeOutput(exported string) error {
	return os.WriteFile(output, []byte(exported), 0o600)
}

// processPaths handles both single file and directory processing.
func processPaths(input string, exclude []string) ([]MsgID, error) {
	stat, err := os.Stat(input)
	if err != nil {
		return nil, err
	}

	if stat.IsDir() {
		return processDirectory(input, exclude)
	}
	return processFile(input)
}

// processDirectory walks through a directory and processes Go files.
func processDirectory(dir string, exclude []string) ([]MsgID, error) {
	var msgids []MsgID

	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if shouldSkipFile(path, info, err, exclude) {
			return nil
		}
		fmt.Println(path)

		ids, err := processFile(path)
		if err != nil {
			return err
		}
		msgids = append(msgids, ids...)
		return nil
	})

	return msgids, err
}

// processFile handles single file processing.
func processFile(path string) ([]MsgID, error) {
	f, err := NewFile(path)
	if err != nil {
		return nil, err
	}

	if !f.HasImport {
		return nil, nil
	}

	return f.MsgIDs(), nil
}

// shouldSkipFile determines if a file should be skipped during processing.
func shouldSkipFile(path string, info fs.FileInfo, err error, exclude []string) bool {
	if err != nil || info.IsDir() {
		return true
	}

	if filepath.Ext(path) != ".go" {
		return true
	}

	return isExcludedPath(path, exclude)
}

// isExcludedPath checks if a path is in the exclude list.
func isExcludedPath(path string, exclude []string) bool {
	return slices.ContainsFunc(exclude, func(s string) bool {
		abs1, err := filepath.Abs(s)
		if err != nil {
			color.Errorln(err)
			return false
		}

		abs2, err := filepath.Abs(path)
		if err != nil {
			color.Errorln(err)
			return false
		}

		return abs1 == abs2
	})
}
