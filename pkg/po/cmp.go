package po

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/Tom5521/xgotext/internal/util"
	fuzzy "github.com/paul-mannino/go-fuzzywuzzy"
)

func EntryMatchRatio(e1, e2 Entry) int {
	ratio := fuzzy.Ratio

	var ratios []int
	if e1.ID != "" || e2.ID != "" {
		ratios = append(ratios, ratio(e1.ID, e2.ID))
	}
	if e1.HasContext() || e2.HasContext() {
		ratios = append(ratios, ratio(e1.Context, e2.Context))
	}
	if e1.IsPlural() {
		ratios = append(ratios, ratio(e1.Plural, e2.Plural))
	}

	return util.Average(ratios...)
}

func CompareEntry(a, b Entry) int {
	obsolete := CompareEntryByObsolete(a, b)
	if obsolete != 0 {
		return obsolete
	}
	fuzzy := CompareEntryByFuzzy(a, b)
	if fuzzy != 0 {
		return fuzzy
	}

	return CompareEntryByLocation(a, b)
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
		return slices.CompareFunc(a.Plurals, b.Plurals, ComparePluralEntryByStr)
	}
	return strings.Compare(a.Str, b.Str)
}

func CompareEntryByLine(a, b Entry) int {
	return slices.CompareFunc(a.Locations, b.Locations, CompareLocationByLine)
}

func CompareEntryByID(a, b Entry) int {
	return strings.Compare(a.ID, b.ID)
}

func CompareEntryByFile(a, b Entry) int {
	return slices.CompareFunc(a.Locations, b.Locations, CompareLocationByFile)
}

func ComparePluralEntryByStr(a, b PluralEntry) int {
	return strings.Compare(a.Str, b.Str)
}

func ComparePluralEntryByID(a, b PluralEntry) int {
	return a.ID - b.ID
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
