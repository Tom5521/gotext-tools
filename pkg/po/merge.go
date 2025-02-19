package po

func MergeFiles(fuzzyMatch bool, base *File, files ...*File) {
	for _, file := range files {
		base.Name += "_" + file.Name
		base.Entries = append(base.Entries, file.Entries...)
	}

	if fuzzyMatch {
		base.Entries = base.Entries.FuzzySolve()
	} else {
		base.Entries = base.Entries.Solve()
	}
	base.Entries = base.Entries.Sort()
}
