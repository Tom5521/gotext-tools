package types

import (
	"path/filepath"
	"slices"
	"strings"
)

type Entries []Entry

func (e Entries) Sort() Entries {
	groupsMap := make(map[string]Entries)
	for _, entry := range e {
		file := ""
		if len(entry.Locations) > 0 {
			file = filepath.Clean(entry.Locations[0].File)
		}
		groupsMap[file] = append(groupsMap[file], entry)
	}

	for k, group := range groupsMap {
		groupsMap[k] = group.SortByLine()
	}

	fileKeys := make([]string, 0, len(groupsMap))
	for file := range groupsMap {
		fileKeys = append(fileKeys, file)
	}
	slices.SortFunc(fileKeys, strings.Compare)

	var sortedEntries Entries
	for _, file := range fileKeys {
		sortedEntries = append(sortedEntries, groupsMap[file]...)
	}

	return sortedEntries
}

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

func (e Entries) SortByID() Entries {
	sorted := make(Entries, len(e))
	copy(sorted, e)
	slices.SortFunc(sorted, func(a, b Entry) int {
		return strings.Compare(a.ID, b.ID)
	})
	return sorted
}

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
