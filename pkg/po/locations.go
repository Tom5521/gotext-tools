package po

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/Tom5521/xgotext/internal/util"
)

// Location represents the location of a translation string in the source code.
type Location struct {
	Line int
	File string
}

func (l Location) Equal(l2 Location) bool {
	return util.Equal(l, l2)
}

type Locations []Location

func (l Locations) Equal(l2 Locations) bool {
	return util.Equal(l, l2)
}

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
