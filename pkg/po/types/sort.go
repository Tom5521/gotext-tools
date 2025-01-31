package types

import (
	"path/filepath"
	"slices"
	"strings"
)

func SortEntriesByID(e []Entry) {
	slices.SortFunc(e, func(a, b Entry) int {
		return strings.Compare(a.ID, b.ID)
	})
}

func SortEntriesByLine(e []Entry) {
	slices.SortFunc(e, func(a, b Entry) int {
		if len(a.Locations) == 0 {
			return 1
		}
		if len(b.Locations) == 0 {
			return -1
		}
		return a.Locations[0].Line - b.Locations[0].Line
	})
}

func SortEntriesByFile(e []Entry) {
	slices.SortFunc(e, func(a, b Entry) int {
		if len(a.Locations) == 0 {
			return 1
		}
		if len(b.Locations) == 0 {
			return -1
		}
		return strings.Compare(a.Locations[0].File, b.Locations[0].File)
	})
}

func SortEntriesByFileAndLine(entries []Entry) {
	groupsMap := make(map[string][]Entry)
	for _, e := range entries {
		file := ""
		if len(e.Locations) > 0 {
			file = filepath.Clean(e.Locations[0].File)
		}
		groupsMap[file] = append(groupsMap[file], e)
	}

	for _, group := range groupsMap {
		slices.SortFunc(group, func(a, b Entry) int {
			if len(a.Locations) == 0 {
				return 1
			}
			if len(b.Locations) == 0 {
				return -1
			}
			return a.Locations[0].Line - b.Locations[0].Line
		})
	}

	fileKeys := make([]string, 0, len(groupsMap))
	for file := range groupsMap {
		fileKeys = append(fileKeys, file)
	}
	slices.SortFunc(fileKeys, strings.Compare)

	var sortedEntries []Entry
	for _, file := range fileKeys {
		sortedEntries = append(sortedEntries, groupsMap[file]...)
	}

	copy(entries, sortedEntries)
}
