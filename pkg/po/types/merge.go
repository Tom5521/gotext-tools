package types

func MergeFiles(base *File, files ...*File) {
	for _, file := range files {
		base.Name += "_" + file.Name
		base.Entries = append(base.Entries, file.Entries...)
	}

	base.Entries = base.Entries.Solve().Sort()
}
