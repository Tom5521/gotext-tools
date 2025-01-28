package types

import "github.com/Tom5521/xgotext/internal/util"

type HeaderField struct {
	Key   string
	Value string
}

type Header struct {
	Values []HeaderField
}

type File struct {
	Name     string
	Header   Header
	Nplurals int
	Entries  []Entry
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
