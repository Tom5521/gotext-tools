package po

import (
	"path/filepath"
	"strings"

	"github.com/Tom5521/gotext-tools/v2/internal/slices"
	"github.com/Tom5521/gotext-tools/v2/internal/util"
	fuzzy "github.com/paul-mannino/go-fuzzywuzzy"
)

func EntryIDMatchRatio(e1, e2 Entry) int {
	ratio := fuzzy.Ratio

	var ratios []int
	if e1.ID != "" || e2.ID != "" {
		ratios = append(ratios, ratio(e1.ID, e2.ID))
	}
	if e1.HasContext() || e2.HasContext() {
		ratios = append(ratios, ratio(e1.Context, e2.Context))
	}
	if e1.IsPlural() && e2.IsPlural() {
		ratios = append(ratios, ratio(e1.Plural, e2.Plural))
	}

	return util.Average(ratios...)
}

func EntryStrMatchRatio(e1, e2 Entry) int {
	ratio := fuzzy.Ratio

	var ratios []int
	if e1.Str != "" || e2.Str != "" {
		ratios = append(ratios, ratio(e1.ID, e2.ID))
	}
	if e1.IsPlural() || e2.IsPlural() {
		for i, pe := range e1.Plurals {
			if i >= 0 && i < len(e2.Plurals) {
				ratios = append(ratios, ratio(pe.Str, e2.Plurals[i].Str))
			}
		}
	}

	return util.Average(ratios...)
}

func CompareEntriesFunc(a, b Entries, cmp Cmp[Entry]) int {
	return slices.CompareFunc(a, b, cmp)
}

func CompareEntry(a, b Entry) int {
	if cmp := CompareEntryByObsolete(a, b); cmp != 0 {
		return cmp
	}

	if cmp := CompareEntryByFuzzy(a, b); cmp != 0 {
		return cmp
	}

	if cmp := CompareEntryByLocation(a, b); cmp != 0 {
		return cmp
	}

	return CompareEntryByID(a, b)
}

func CompareEntryByObsolete(a, b Entry) int {
	if !a.Obsolete && b.Obsolete {
		return -1
	} else if a.Obsolete && !b.Obsolete {
		return 1
	}

	return 0
}

func CompareEntryByFuzzy(a, b Entry) int {
	aContains := a.IsFuzzy()
	bContains := b.IsFuzzy()

	if !aContains && bContains {
		return -1
	} else if aContains && !bContains {
		return 1
	}

	return 0
}

func CompareEntryByLocation(a, b Entry) int {
	return slices.CompareFunc(a.Locations, b.Locations, CompareLocation)
}

func CompareEntryByStr(a, b Entry) int {
	if a.IsPlural() && b.IsPlural() {
		if i := ComparePluralEntriesFunc(a.Plurals, b.Plurals, ComparePluralEntryByStr); i != 0 {
			return i
		}
	}

	return strings.Compare(a.Str, b.Str)
}

func CompareEntryByLine(a, b Entry) int {
	return slices.CompareFunc(a.Locations, b.Locations, CompareLocationByLine)
}

func CompareEntryByID(a, b Entry) int {
	return strings.Compare(a.UnifiedID(), b.UnifiedID())
}

func CompareEntryByFile(a, b Entry) int {
	return slices.CompareFunc(a.Locations, b.Locations, CompareLocationByFile)
}

func ComparePluralEntriesFunc(a, b PluralEntries, cmp Cmp[PluralEntry]) int {
	return slices.CompareFunc(a, b, cmp)
}

func ComparePluralEntry(a, b PluralEntry) int {
	id := ComparePluralEntryByID(a, b)
	if id != 0 {
		return id
	}
	return ComparePluralEntryByStr(a, b)
}

func ComparePluralEntryByStr(a, b PluralEntry) int {
	return strings.Compare(a.Str, b.Str)
}

func ComparePluralEntryByID(a, b PluralEntry) int {
	return a.ID - b.ID
}

func CompareLocationsFunc(a, b Locations, cmp Cmp[Location]) int {
	return slices.CompareFunc(a, b, cmp)
}

func CompareLocation(a, b Location) int {
	if file := CompareLocationByFile(a, b); file != 0 {
		return file
	}

	return CompareLocationByLine(a, b)
}

func CompareLocationByLine(a, b Location) int {
	return a.Line - b.Line
}

func CompareLocationByFile(a, b Location) int {
	return strings.Compare(filepath.Clean(a.File), filepath.Clean(b.File))
}
