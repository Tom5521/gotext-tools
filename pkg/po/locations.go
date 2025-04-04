package po

import (
	"slices"

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
	slices.SortFunc(l, CompareLocation)
	return l
}

func (l Locations) SortByLine() Locations {
	slices.SortFunc(l, CompareLocationByLine)

	return l
}

func (l Locations) SortByFile() Locations {
	slices.SortFunc(l, CompareLocationByFile)

	return l
}
