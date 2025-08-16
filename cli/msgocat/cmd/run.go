package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Tom5521/gotext-tools/v2/pkg/po"
	"github.com/Tom5521/gotext-tools/v2/pkg/po/parse"
	"github.com/spf13/cobra"
)

const (
	fileTemplate = "#-#-#-#-#  %s  #-#-#-#-#"
)

// entryInfo tracks metadata about merged entries.
type entryInfo struct {
	index int    // Index in the entries slice
	file  string // First file where the entry was found
	ncols uint   // Number of collisions detected
}

// run is the main execution function for the command.
func run(cmd *cobra.Command, args []string) error {
	args, err := expandArgsWithFiles(args)
	if err != nil {
		return err
	}

	var (
		entries    po.Entries
		entriesMap = make(map[string]entryInfo)
	)

	for _, file := range args {
		if err := processPOFile(file, &entries, entriesMap, len(args)); err != nil {
			return err
		}
	}
	return nil
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

// processPOFile processes a single PO file and merges its entries.
func processPOFile(
	filename string,
	entries *po.Entries,
	entriesMap map[string]entryInfo,
	totalFiles int,
) error {
	poFile, err := parsePOFile(filename)
	if err != nil {
		return err
	}

	for i := range poFile.Entries {
		if err := processEntry(&poFile.Entries[i], poFile.Name, entries, entriesMap, totalFiles); err != nil {
			return err
		}
	}
	return nil
}

// parsePOFile parses a PO file from disk.
func parsePOFile(filename string) (*po.File, error) {
	fileRef, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer fileRef.Close()

	poFile, err := parse.PoFromFile(fileRef, parse.PoWithCleanDuplicates(false))
	if err != nil {
		return nil, err
	}

	if errs := poFile.Validate(); len(errs) > 0 {
		return nil, errs[0]
	}
	return poFile, nil
}

// processEntry handles an individual PO entry including merging logic.
func processEntry(
	entry *po.Entry,
	filename string,
	entries *po.Entries,
	entriesMap map[string]entryInfo,
	totalFiles int,
) error {
	if !entry.HasComments() || totalFiles <= 1 {
		return nil
	}

	// Add file header to comments
	entry.Comments = append([]string{fmt.Sprintf(fileTemplate, filename)}, entry.Comments...)

	id := entry.UnifiedID()
	info, exists := entriesMap[id]

	if !exists {
		// Add new entry
		*entries = append(*entries, *entry)
		entriesMap[id] = entryInfo{
			index: len(*entries) - 1,
			file:  filename,
			ncols: 0,
		}
		return nil
	}

	existingEntry := &(*entries)[info.index]
	if existingEntry.UnifiedStr() == entry.UnifiedStr() {
		return nil // Duplicate translation
	}

	// Handle collision
	info.ncols++
	entriesMap[id] = info
	return mergeEntry(existingEntry, entry, info.file, filename, info.ncols)
}

// mergeEntry merges two conflicting PO entries.
func mergeEntry(
	existing *po.Entry,
	new *po.Entry,
	existingFile, newFile string,
	collisionCount uint,
) error {
	existing.Flags = append(existing.Flags, "fuzzy")

	// Merge main string
	existing.Str = processString(
		existingFile,
		newFile,
		existing.Str,
		new.Str,
		collisionCount == 1,
	)

	// Merge plural forms if both entries are plural
	if existing.IsPlural() && new.IsPlural() {
		mergePlurals(existing, new, existingFile)
	}
	return nil
}

// mergePlurals handles merging of plural form strings.
func mergePlurals(existing *po.Entry, new *po.Entry, existingFile string) {
	for i := range existing.Plurals {
		if i >= len(new.Plurals) {
			// Add origin marker if new entry is missing this plural form
			existing.Plurals[i].Str = addFileHeader(existingFile, existing.Plurals[i].Str)
			continue
		}

		// Concatenate translations
		existing.Plurals[i].Str = strings.Join([]string{
			existing.Plurals[i].Str,
			new.Plurals[i].Str,
		}, "\n")
	}
}

// processString merges two translation strings with file markers.
func processString(f1, f2, str1, str2 string, isNew bool) string {
	if str1 == str2 {
		return str1
	}

	lines := strings.Split(str1, "\n")
	if isNew {
		lines = append([]string{fmt.Sprintf(fileTemplate, f1)}, lines...)
	}
	return strings.Join(append(lines, fmt.Sprintf(fileTemplate, f2), "\n"), "\n")
}

// addFileHeader prepends a file marker to a string.
func addFileHeader(file, str string) string {
	return strings.Join(append(
		[]string{fmt.Sprintf(fileTemplate, file)},
		strings.Split(str, "\n")...,
	), "\n")
}
