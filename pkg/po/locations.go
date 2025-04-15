package po

import (
	"github.com/Tom5521/gotext-tools/internal/slices"

	"github.com/Tom5521/gotext-tools/internal/util"
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

func (l Locations) IsSorted() bool {
	return l.IsSortedFunc(CompareLocation)
}

func (l Locations) IsSortedFunc(cmp Cmp[Location]) bool {
	return slices.IsSortedFunc(l, cmp)
}

func (l Locations) Sort() Locations {
	slices.SortFunc(l, CompareLocation)
	return l
}

func (l Locations) SortFunc(cmp Cmp[Location]) Locations {
	slices.SortFunc(l, cmp)
	return l
}
