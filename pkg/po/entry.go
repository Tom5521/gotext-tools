package po

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/Tom5521/xgotext/internal/util"
)

// Location represents the location of a translation string in the source code.
type Location struct {
	Line int    // The line number of the translation.
	File string // The file name where the translation is located.
}

type Locations []Location

func (l Locations) Sort() Locations {
	groups := make(map[string]Locations)
	for _, l2 := range l {
		l2.File = filepath.Clean(l2.File)
		groups[l2.File] = append(groups[l2.File], l2)
	}

	for k, l2 := range groups {
		groups[k] = l2.SortByLine()
	}

	fileKeys := make([]string, 0, len(groups))
	for k := range groups {
		fileKeys = append(fileKeys, k)
	}
	slices.Sort(fileKeys)

	var sorted Locations
	for _, file := range fileKeys {
		sorted = append(sorted, groups[file]...)
	}

	return sorted
}

func (l Locations) SortByLine() Locations {
	slices.SortFunc(l, func(a, b Location) int {
		return a.Line - b.Line
	})

	return l
}

func (l Locations) SortByFile() Locations {
	slices.SortFunc(l, func(a, b Location) int {
		return strings.Compare(a.File, b.File)
	})

	return l
}

type PluralEntries []PluralEntry

func (p PluralEntries) Sort() PluralEntries {
	slices.SortFunc(p, func(a, b PluralEntry) int {
		return a.ID - b.ID
	})

	return p
}

type PluralEntry struct {
	ID  int
	Str string
}

// Entry represents a translatable string, including its context, plural forms,
// and source code locations.
type Entry struct {
	// Comments.

	Flags             []string
	Comments          []string
	ExtractedComments []string
	Previous          []string

	// Fields.

	ID        string // The original string to be translated.
	Context   string // The context in which the string is used (optional).
	Plural    string // The plural form of the string (optional).
	Plurals   PluralEntries
	Str       string
	Locations Locations // A list of source code locations for the string.
}

func (e Entry) String() string {
	return util.Format(e)
}
