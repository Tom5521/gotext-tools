package po

import (
	"path/filepath"
	"strings"

	"github.com/Tom5521/gotext-tools/v2/internal/slices"
	"github.com/Tom5521/gotext-tools/v2/internal/util"
	fuzzy "github.com/paul-mannino/go-fuzzywuzzy"
)

// EntryIDMatchRatio returns the average fuzzy match ratio between the identifiers of two entries.
// It compares the ID, context (if present), and plural ID (if both are plural entries).
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

// EntryStrMatchRatio returns the average fuzzy match ratio between the translation strings
// of two entries, including plural forms if applicable.
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

// CompareEntriesFunc compares two slices of entries using the given comparison function.
func CompareEntriesFunc(a, b Entries, cmp Cmp[Entry]) int {
	return slices.CompareFunc(a, b, cmp)
}

// CompareEntry compares two entries based on their obsolescence, fuzzy flag,
// location, and ID (in that order).
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

// CompareEntryByObsolete compares entries by their obsolescence status.
func CompareEntryByObsolete(a, b Entry) int {
	if !a.Obsolete && b.Obsolete {
		return -1
	} else if a.Obsolete && !b.Obsolete {
		return 1
	}
	return 0
}

// CompareEntryByFuzzy compares entries by their fuzzy flag.
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

// CompareEntryByLocation compares entries by their location slices.
func CompareEntryByLocation(a, b Entry) int {
	return slices.CompareFunc(a.Locations, b.Locations, CompareLocation)
}

// CompareEntryByStr compares entries by their translation strings, including plural forms.
func CompareEntryByStr(a, b Entry) int {
	if a.IsPlural() && b.IsPlural() {
		if i := ComparePluralEntriesFunc(a.Plurals, b.Plurals, ComparePluralEntryByStr); i != 0 {
			return i
		}
	}
	return strings.Compare(a.Str, b.Str)
}

// CompareEntryByLine compares entries by line numbers in their locations.
func CompareEntryByLine(a, b Entry) int {
	return slices.CompareFunc(a.Locations, b.Locations, CompareLocationByLine)
}

// CompareEntryByID compares entries by their unified ID.
func CompareEntryByID(a, b Entry) int {
	return strings.Compare(a.UnifiedID(), b.UnifiedID())
}

// CompareEntryByFile compares entries by file paths in their locations.
func CompareEntryByFile(a, b Entry) int {
	return slices.CompareFunc(a.Locations, b.Locations, CompareLocationByFile)
}

// ComparePluralEntriesFunc compares two plural entry slices using the provided comparison function.
func ComparePluralEntriesFunc(a, b PluralEntries, cmp Cmp[PluralEntry]) int {
	return slices.CompareFunc(a, b, cmp)
}

// ComparePluralEntry compares two plural entries by ID and then by translation string.
func ComparePluralEntry(a, b PluralEntry) int {
	if id := ComparePluralEntryByID(a, b); id != 0 {
		return id
	}
	return ComparePluralEntryByStr(a, b)
}

// ComparePluralEntryByStr compares two plural entries by their translation strings.
func ComparePluralEntryByStr(a, b PluralEntry) int {
	return strings.Compare(a.Str, b.Str)
}

// ComparePluralEntryByID compares two plural entries by their numeric IDs.
func ComparePluralEntryByID(a, b PluralEntry) int {
	return a.ID - b.ID
}

// CompareLocationsFunc compares two slices of locations using the given comparison function.
func CompareLocationsFunc(a, b Locations, cmp Cmp[Location]) int {
	return slices.CompareFunc(a, b, cmp)
}

// CompareLocation compares two locations by file path and then by line number.
func CompareLocation(a, b Location) int {
	if file := CompareLocationByFile(a, b); file != 0 {
		return file
	}
	return CompareLocationByLine(a, b)
}

// CompareLocationByLine compares two locations by line number.
func CompareLocationByLine(a, b Location) int {
	return a.Line - b.Line
}

// CompareLocationByFile compares two locations by cleaned file path string.
func CompareLocationByFile(a, b Location) int {
	return strings.Compare(filepath.Clean(a.File), filepath.Clean(b.File))
}
