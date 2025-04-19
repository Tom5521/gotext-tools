package po

import (
	"github.com/Tom5521/gotext-tools/v2/internal/util"
)

type File struct {
	Name string
	Entries
}

func NewFile(name string, entries ...Entry) *File {
	return &File{name, entries}
}

func (f File) Equal(f2 File) bool {
	return util.Equal(f, f2)
}

func (f *File) Set(id, context string, e Entry) {
	index := f.Index(id, context)
	if index == -1 {
		f.Entries = append(f.Entries, e)
		return
	}
	f.Entries[index] = e
}

func (f File) LoadByUnifiedID(uid string) string {
	i := f.IndexByUnifiedID(uid)
	if i == -1 {
		return ""
	}
	return f.Entries[i].Str
}

func (f File) Load(id string, context string) string {
	i := f.Index(id, context)
	if i == -1 {
		return ""
	}

	return f.Entries[i].Str
}
