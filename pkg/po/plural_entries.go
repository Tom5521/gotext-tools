package po

import (
	"slices"

	"github.com/Tom5521/xgotext/internal/util"
)

type PluralEntry struct {
	ID  int
	Str string
}

func (p PluralEntry) Equal(p2 PluralEntry) bool {
	return util.Equal(p, p2)
}

type PluralEntries []PluralEntry

func (p PluralEntries) Equal(p2 PluralEntries) bool {
	return util.Equal(p2, p)
}

func (p PluralEntries) Sort() PluralEntries {
	slices.SortFunc(p, func(a, b PluralEntry) int {
		return a.ID - b.ID
	})

	return p
}
