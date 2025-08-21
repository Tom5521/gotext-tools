package cmd

import (
	"fmt"
	"strings"

	"github.com/Tom5521/gotext-tools/v2/internal/slices"
	"github.com/Tom5521/gotext-tools/v2/pkg/po"
)

// TODO: Implement this on the main library.

// By now, I'll implement this with functional programming.

const (
	fileTemplate = "#-#-#-#-#  %s  #-#-#-#-#"
)

type processState struct {
	entriesMap   map[string]entryInfo
	entriesSlice po.Entries
	// totalFiles   int
}

// entryInfo tracks metadata about merged entries.
type entryInfo struct {
	index      int    // Index in the entries slice
	file       string // First file where the entry was found
	timesFound int
}

// mergePoFiles processes a single PO file and merges its entries.
func mergePoFiles(
	filename string,
	state *processState,
) error {
	poFile, err := parsePOFile(filename)
	if err != nil {
		return err
	}

	for _, entry := range poFile.Entries {
		handleEntry(entry, poFile.Name, state)
	}
	return nil
}

// handleEntry handles an individual PO entry including merging logic.
func handleEntry(
	entry po.Entry,
	filename string,
	state *processState,
) {
	uid := entry.UnifiedID()
	original, duplicated := state.entriesMap[uid]
	if duplicated {
		state.entriesSlice[original.index] = mergeEntry(
			state.entriesSlice[original.index],
			entry,
			original.file,
			filename,
		)
		original.timesFound++
		state.entriesMap[uid] = original
		return
	}

	state.entriesMap[uid] = entryInfo{
		index:      len(state.entriesSlice),
		file:       filename,
		timesFound: 1,
	}
	state.entriesSlice = append(state.entriesSlice, entry)
}

// mergeEntry merges two conflicting PO entries.
func mergeEntry(
	originalEntry po.Entry,
	newEntry po.Entry,
	originalFilename, newFilename string,
) po.Entry {
	mergedEntry := po.Entry{
		Flags: mergeComments(
			originalEntry.Flags,
			newEntry.Flags,
			originalFilename,
			newFilename,
		),
		Comments: mergeComments(
			originalEntry.Comments,
			newEntry.Comments,
			originalFilename,
			newFilename,
		),
		ExtractedComments: mergeComments(
			originalEntry.ExtractedComments,
			newEntry.ExtractedComments,
			originalFilename,
			newFilename,
		),
		Previous: mergeComments(
			originalEntry.Previous,
			newEntry.Previous,
			originalFilename,
			newFilename,
		),

		ID:      originalEntry.ID,
		Context: originalEntry.Context,
		Plural:  originalEntry.Plural,
		Plurals: mergePlurals(
			originalEntry,
			newEntry,
			originalFilename,
			newFilename,
		),
		Str: mergeStrings(
			originalEntry.Str,
			newEntry.Str,
			originalFilename,
			newFilename,
		),
		Locations: append(
			originalEntry.Locations,
			newEntry.Locations...,
		),
	}

	if !mergedEntry.IsFuzzy() &&
		(originalEntry.FullUnifiedID() != newEntry.FullUnifiedID() ||
			originalEntry.UnifiedStr() != newEntry.UnifiedStr()) {
		mergedEntry.Flags = append(mergedEntry.Flags, "fuzzy")
	}

	return mergedEntry
}

// mergePlurals handles merging of plural form strings.
func mergePlurals(
	originalEntry, newEntry po.Entry,
	originalFilename, newFilename string,
) po.PluralEntries {
	if !originalEntry.IsPlural() {
		return originalEntry.Plurals
	}

	mergedEntries := make(po.PluralEntries, 0, len(originalEntry.Plurals)+len(newEntry.Plurals))

	for i, original := range originalEntry.Plurals {
		if i >= len(newEntry.Plurals) {
			mergedEntries = append(mergedEntries, original)
			continue
		}

		newPlural := newEntry.Plurals[i]
		mergedEntries = append(mergedEntries, po.PluralEntry{
			ID: original.ID,
			Str: mergeStrings(
				original.Str, newPlural.Str,
				originalFilename, newFilename,
			),
		})
	}

	return mergedEntries
}

func mergeStrings(originalStr, newStr string, originalFilename, newFilename string) string {
	if originalStr == newStr {
		return originalStr
	}
	if originalStr == "" {
		return newStr
	}
	if newStr == "" {
		return originalStr
	}

	var builder strings.Builder
	fmt.Fprintf(&builder, fileTemplate+"\n", originalFilename)
	fmt.Fprintln(&builder, originalStr)
	fmt.Fprintf(&builder, fileTemplate+"\n", newFilename)
	fmt.Fprint(&builder, newStr)
	return builder.String()
}

func mergeComments(
	originalComments, newComments []string,
	originalFilename, newFilename string,
) []string {
	if slices.Equal(originalComments, newComments) {
		return originalComments
	}

	if len(originalComments) == 0 {
		return append(
			[]string{fmt.Sprintf(fileTemplate, newFilename)}, newComments...,
		)
	}
	if len(newComments) == 0 {
		return append(
			[]string{fmt.Sprintf(fileTemplate, originalFilename)},
			originalComments...,
		)
	}

	return append(
		append(
			[]string{fmt.Sprintf(fileTemplate, originalFilename)},
			originalComments...,
		),
		append(
			[]string{fmt.Sprintf(fileTemplate, newFilename)},
			newComments...,
		)...,
	)
}
