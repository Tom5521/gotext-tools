package util

import (
	"log"
	"path/filepath"

	"github.com/Tom5521/gotext-tools/v2/internal/slices"
	krfs "github.com/kr/fs"
)

// ShouldSkipFile determines if a file should be skipped during processing.
func ShouldSkipFile(
	w *krfs.Walker,
	excludedPaths []string,
	seenMap *map[string]bool,
	logger *log.Logger,
) bool {
	if w.Err() != nil || w.Stat().IsDir() {
		return true
	}

	if filepath.Ext(w.Path()) != ".go" {
		return true
	}

	abs, err := filepath.Abs(w.Path())
	if err != nil {
		return true
	}

	_, seen := (*seenMap)[abs]
	if seen {
		logger.Printf("skipping duplicated file: %s\n", w.Path())
		return true
	}
	(*seenMap)[abs] = true

	return isExcludedPath(w.Path(), excludedPaths)
}

// isExcludedPath checks if a path is in the exclude list defined in the configuration.
func isExcludedPath(path string, exclude []string) bool {
	return slices.ContainsFunc(exclude, func(s string) bool {
		abs1, err1 := filepath.Abs(s)
		abs2, err2 := filepath.Abs(path)
		return (abs1 == abs2) && (err1 == nil && err2 == nil)
	})
}
