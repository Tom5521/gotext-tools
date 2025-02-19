package po

import (
	"path/filepath"
	"slices"
	"strings"

	fuzzy "github.com/paul-mannino/go-fuzzywuzzy"
)

// Entries represents a collection of Entry objects.
type Entries []Entry

func (e Entries) Index(id, context string) int {
	for i, entry := range e {
		if entry.ID == id && entry.Context == context {
			return i
		}
	}

	return -1
}

// Sort organizes the entries by grouping them by file and sorting them by line.
func (e Entries) Sort() Entries {
	groupsMap := make(map[string]Entries)

	// Group entries by file.
	for _, entry := range e {
		var file string
		if len(entry.Locations) > 0 {
			file = filepath.Clean(entry.Locations[0].File)
		}
		groupsMap[file] = append(groupsMap[file], entry)
	}

	// Sort each group by line.
	for k, group := range groupsMap {
		groupsMap[k] = group.SortByLine()
	}

	// Get sorted file names.
	fileKeys := make([]string, 0, len(groupsMap))
	for file := range groupsMap {
		fileKeys = append(fileKeys, file)
	}
	slices.SortFunc(fileKeys, strings.Compare)

	// Concatenate groups into a single sorted list.
	var sortedEntries Entries
	for _, file := range fileKeys {
		sortedEntries = append(sortedEntries, groupsMap[file]...)
	}

	return sortedEntries
}

// SortByFile sorts the entries by the file name of the first location.
func (e Entries) SortByFile() Entries {
	sorted := make(Entries, len(e))
	copy(sorted, e)
	slices.SortFunc(sorted, func(a, b Entry) int {
		if len(a.Locations) == 0 {
			return 1
		}
		if len(b.Locations) == 0 {
			return -1
		}
		return strings.Compare(a.Locations[0].File, b.Locations[0].File)
	})
	return sorted
}

// SortByID sorts the entries by their ID.
func (e Entries) SortByID() Entries {
	sorted := make(Entries, len(e))
	copy(sorted, e)
	slices.SortFunc(sorted, func(a, b Entry) int {
		return strings.Compare(a.ID, b.ID)
	})
	return sorted
}

// SortByLine sorts the entries by line number in their first location.
func (e Entries) SortByLine() Entries {
	sorted := make(Entries, len(e))
	copy(sorted, e)
	slices.SortFunc(sorted, func(a, b Entry) int {
		if len(a.Locations) == 0 {
			return 1
		}
		if len(b.Locations) == 0 {
			return -1
		}
		return a.Locations[0].Line - b.Locations[0].Line
	})
	return sorted
}

// CleanDuplicates removes duplicate entries with the same ID and context, merging their locations.
func (e Entries) CleanDuplicates() Entries {
	var cleaned Entries
	seenID := make(map[string]int)

	for _, translation := range e {
		idIndex, ok := seenID[translation.ID]
		if ok {
			if translation.Context == cleaned[idIndex].Context {
				cleaned[idIndex].Locations = append(
					cleaned[idIndex].Locations,
					translation.Locations...)
				continue
			}
		}
		seenID[translation.ID] = len(cleaned)
		cleaned = append(cleaned, translation)
	}

	return cleaned
}

// Solve processes a list of translation entries and merges those with the same ID and context,
// keeping the most complete translation. If two entries have the same ID and context, the one
// with a non-empty translation string is retained. Additionally, if the entries are similar but not
// identical, the resulting entry is marked as "fuzzy". The locations of the merged entries are combined.
func (e Entries) Solve() Entries {
	var cleaned Entries
	seenID := make(map[string]int)

	for _, translation := range e {
		idIndex, ok := seenID[translation.ID]
		if ok {
			if translation.Context == cleaned[idIndex].Context {
				// If the new entry has a translation and the previous one does not, replace it.
				if translation.Str != "" && cleaned[idIndex].Str == "" {
					cleaned[idIndex] = translation
				}

				// Combine the locations of the merged entries.
				cleaned[idIndex].Locations = append(
					cleaned[idIndex].Locations,
					translation.Locations...)
				continue
			}
		}
		seenID[translation.ID] = len(cleaned)
		cleaned = append(cleaned, translation)
	}

	return cleaned
}

func (e Entries) FuzzyFind(id, context string) int {
	for i, entry := range e {
		if fuzzy.Ratio(entry.ID, id) >= 80 && entry.Context == context {
			return i
		}
	}

	return -1
}

func (e Entries) FuzzySolve() (cleaned Entries) {
	var dupedGroups []Entries

	find := func(e Entry) int {
		for i, group := range dupedGroups {
			if len(group) > 0 {
				if fuzzy.Ratio(group[0].ID, e.ID) >= 80 &&
					group[0].Context == e.Context {
					return i
				}
			}
		}
		return -1
	}

	// Collect duplicates
	for _, entry := range e {
		groupIndex := find(entry)
		if groupIndex == -1 {
			dupedGroups = append(dupedGroups, []Entry{entry})
		} else {
			dupedGroups[groupIndex] = append(dupedGroups[groupIndex], entry)
		}
	}
	// Clean duplicates
	for _, group := range dupedGroups {
		if len(group) == 1 {
			cleaned = append(cleaned, group[0])
			continue
		}
		entry := group.Solve()[0]
		if !slices.Contains(entry.Flags, "fuzzy") {
			entry.Flags = append(entry.Flags, "fuzzy")
		}
		cleaned = append(cleaned, entry)
	}

	return
}
