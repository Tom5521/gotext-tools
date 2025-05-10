package po

import (
	"github.com/Tom5521/gotext-tools/v2/internal/util"
)

// File represents a PO file with a name and a list of entries.
type File struct {
	Name    string // File name or path.
	Entries        // List of translation entries.
}

// NewFile returns a new *File with the given name and entries.
func NewFile(name string, entries ...Entry) *File {
	return &File{name, entries}
}

// Equal reports whether f and f2 contain the same entries and metadata.
func (f File) Equal(f2 File) bool {
	return util.Equal(f, f2)
}

// Set sets the entry with the given id and context to e.
// If the entry exists, it is replaced; otherwise, it is appended.
func (f *File) Set(id, context string, e Entry) {
	index := f.Index(id, context)
	if index == -1 {
		f.Entries = append(f.Entries, e)
		return
	}
	f.Entries[index] = e
}

// LoadByUnifiedID returns the translation string for the entry with the given unified ID.
// If no such entry exists, it returns an empty string.
func (f File) LoadByUnifiedID(uid string) string {
	i := f.IndexByUnifiedID(uid)
	if i == -1 {
		return ""
	}
	return f.Entries[i].Str
}

// Load returns the translation string for the entry with the given ID and context.
// If no such entry exists, it returns an empty string.
func (f File) Load(id string, context string) string {
	i := f.Index(id, context)
	if i == -1 {
		return ""
	}
	return f.Entries[i].Str
}

// String returns the formatted representation of the file contents.
func (f File) String() string {
	return util.Format(f)
}
