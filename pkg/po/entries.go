package po

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/Tom5521/xgotext/internal/util"
	fuzzy "github.com/paul-mannino/go-fuzzywuzzy"
)

// Entries represents a collection of Entry objects.
type Entries []Entry

func (e Entries) Equal(e2 Entries) bool {
	return util.Equal(e, e2)
}

func (e Entries) IndexByUnifiedID(uid string) int {
	return slices.IndexFunc(e, func(e Entry) bool {
		return e.UnifiedID() == uid
	})
}

func (e Entries) Index(id, context string) int {
	return slices.IndexFunc(e,
		func(e Entry) bool {
			return e.ID == id && e.Context == context
		},
	)
}

func (e Entries) IsSorted() bool {
	return util.Equal(e, e.Sort())
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
	slices.Sort(fileKeys)

	// Concatenate groups into a single sorted list.
	var sortedEntries Entries
	for _, file := range fileKeys {
		sortedEntries = append(sortedEntries, groupsMap[file]...)
	}

	return sortedEntries.SortByFuzzy()
}

func (e Entries) IsSortedByFuzzy() bool {
	return util.Equal(e, e.SortByFuzzy())
}

func (e Entries) SortByFuzzy() Entries {
	slices.SortFunc(e, func(a, b Entry) int {
		aContains := a.IsFuzzy()
		bContains := b.IsFuzzy()

		switch {
		case aContains == bContains:
			return 0
		case aContains && !bContains:
			return 1
		default:
			return -1
		}
	})

	return e
}

func (e Entries) IsSortedByFile() bool {
	return util.Equal(e, e.SortByFile())
}

// SortByFile sorts the entries by the file name of the first location.
func (e Entries) SortByFile() Entries {
	slices.SortFunc(e, func(a, b Entry) int {
		if len(a.Locations) == 0 {
			return 1
		}
		if len(b.Locations) == 0 {
			return -1
		}
		return strings.Compare(a.Locations[0].File, b.Locations[0].File)
	})
	return e
}

func (e Entries) IsSortedByID() bool {
	return util.Equal(e, e.SortByID())
}

// SortByID sorts the entries by their ID.
func (e Entries) SortByID() Entries {
	slices.SortFunc(e, func(a, b Entry) int {
		return strings.Compare(a.ID, b.ID)
	})
	return e
}

func (e Entries) IsSortedByLine() bool {
	return util.Equal(e, e.SortByLine())
}

// SortByLine sorts the entries by line number in their first location.
func (e Entries) SortByLine() Entries {
	slices.SortFunc(e, func(a, b Entry) int {
		if len(a.Locations) == 0 {
			return 1
		}
		if len(b.Locations) == 0 {
			return -1
		}
		return a.Locations[0].Line - b.Locations[0].Line
	})
	return e
}

func (e Entries) HasDuplicates() bool {
	seen := make(map[string]bool)

	for _, entry := range e {
		uid := entry.UnifiedID()
		_, seened := seen[uid]
		if seened {
			return true
		}

		seen[uid] = true
	}

	return false
}

// CleanDuplicates removes duplicate entries with the same ID and context, merging their locations.
func (e Entries) CleanDuplicates() Entries {
	var cleaned Entries
	seenID := make(map[string]int)

	for _, entry := range e {
		uid := entry.UnifiedID()
		idIndex, ok := seenID[uid]
		if ok {
			cleaned[idIndex].Locations = append(cleaned[idIndex].Locations, entry.Locations...)
			continue
		}
		seenID[uid] = len(cleaned)
		cleaned = append(cleaned, entry)
	}

	return cleaned
}

// Solve processes a list of translation entries and merges those with the same ID and context,
// keeping the most complete translation. If two entries have the same ID and context, the one
// with a non-empty translation string is retained. Additionally, if the entries are similar but not.
func (e Entries) Solve() Entries {
	var cleaned Entries
	seenID := make(map[string]int)

	for _, entry := range e {
		uid := entry.UnifiedID()
		idIndex, ok := seenID[uid]
		if ok {
			// If the new entry has a translation and the previous one does not, replace it.
			if entry.IsPlural() {
				if len(entry.Plurals) != 0 && len(cleaned[idIndex].Plurals) > 0 {
					cleaned[idIndex].Plurals = append(
						entry.Plurals,
						cleaned[idIndex].Plurals...).Solve()
				}
			} else if entry.Str != "" && cleaned[idIndex].Str == "" {
				cleaned[idIndex] = entry
			}

			// Combine the locations of the merged entries.
			cleaned[idIndex].Locations = append(
				cleaned[idIndex].Locations,
				entry.Locations...)
			continue
		}
		seenID[uid] = len(cleaned)
		cleaned = append(cleaned, entry)
	}

	return cleaned
}

func (e Entries) CleanFuzzy() Entries {
	e = slices.DeleteFunc(e, func(e Entry) bool {
		return e.IsFuzzy()
	})
	return e
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
		if !entry.IsFuzzy() {
			entry.Flags = append(entry.Flags, "fuzzy")
		}
		cleaned = append(cleaned, entry)
	}

	return
}
