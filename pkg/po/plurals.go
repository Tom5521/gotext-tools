package po

import (
	"slices"
	"strconv"

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

func (p PluralEntries) Solve() PluralEntries {
	seen := make(map[string]bool)
	var cleaned PluralEntries

	for _, pe := range p {
		id := pe.Str + "\x00" + strconv.Itoa(pe.ID)
		_, seened := seen[id]
		if seened {
			continue
		}
		cleaned = append(cleaned, pe)
	}

	return cleaned
}

func (p PluralEntries) Sort() PluralEntries {
	slices.SortFunc(p, ComparePluralEntryByID)

	return p
}
