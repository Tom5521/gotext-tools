package types

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/Tom5521/xgotext/internal/util"
)

type Entries []Entry

func (e Entries) Sort() Entries {
	groupsMap := make(map[string][]Entry)
	for _, entry := range e {
		file := ""
		if len(entry.Locations) > 0 {
			file = filepath.Clean(entry.Locations[0].File)
		}
		groupsMap[file] = append(groupsMap[file], entry)
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

type File struct {
	Name     string
	Header   Header
	Nplurals int
	Entries  Entries
}

func (f File) LoadID(id string) string {
	for _, entry := range f.Entries {
		if entry.ID == id {
			return entry.Str
		}
	}

	return ""
}

// Location represents the location of a translation string in the source code.
type Location struct {
	Line int    // The line number of the translation.
	File string // The file name where the translation is located.
}

type PluralEntry struct {
	ID  int
	Str string
}

// Entry represents a translatable string, including its context, plural forms,
// and source code locations.
type Entry struct {
	Flags     []string
	ID        string // The original string to be translated.
	Context   string // The context in which the string is used (optional).
	Plural    string // The plural form of the string (optional).
	Plurals   []PluralEntry
	Str       string
	Locations []Location // A list of source code locations for the string.
}

func (e Entry) String() string {
	return util.Format(e)
}
