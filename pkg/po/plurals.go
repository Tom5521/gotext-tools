package po

import (
	"fmt"

	"github.com/Tom5521/gotext-tools/v2/internal/slices"
	"github.com/Tom5521/gotext-tools/v2/internal/util"
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

func (p PluralEntries) Solve() PluralEntries {
	seen := make(map[string]struct{}, len(p))
	cleaned := make(PluralEntries, 0, len(p))

	for _, pe := range p {
		id := fmt.Sprintf("%s\x00%d", pe.Str, pe.ID)
		if _, exists := seen[id]; !exists {
			seen[id] = struct{}{}
			cleaned = append(cleaned, pe)
		}
	}

	return cleaned
}

func (p PluralEntries) IsSorted() bool {
	return slices.IsSortedFunc(p, ComparePluralEntry)
}

func (p PluralEntries) IsSortedFunc(cmp Cmp[PluralEntry]) bool {
	return slices.IsSortedFunc(p, cmp)
}

func (p PluralEntries) Sort() PluralEntries {
	slices.SortFunc(p, ComparePluralEntry)

	return p
}

func (p PluralEntries) SortFunc(cmp Cmp[PluralEntry]) PluralEntries {
	slices.SortFunc(p, cmp)
	return p
}
