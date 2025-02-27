package po

type File struct {
	Entries Entries
	Name    string
}

func NewFile(name string, entries ...Entry) *File {
	f := &File{Name: name, Entries: entries}

	return f
}

func (f File) Header() Header {
	return f.Entries.Header()
}

func (f *File) Set(id, context string, e Entry) {
	index := f.Entries.Index(id, context)
	if index == -1 {
		f.Entries = append(f.Entries, e)
		return
	}
	f.Entries[index] = e
}

func (f File) LoadID(id string, context string) string {
	i := f.Entries.Index(id, context)
	if i == -1 {
		return ""
	}

	return f.Entries[i].Str
}
