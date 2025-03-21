package po

import "strings"

func MergeFiles(fuzzyMatch bool, base *File, files ...*File) {
	names := []string{base.Name}
	for _, file := range files {
		names = append(names, file.Name)
		base.Entries = append(base.Entries, file.Entries...)
	}
	base.Name = strings.Join(names, "_")

	if fuzzyMatch {
		base.Entries = base.Entries.FuzzySolve()
	} else {
		base.Entries = base.Entries.Solve()
	}
	base.Entries = base.Entries.Sort()
}
