package po

import (
	"github.com/Tom5521/gotext-tools/internal/slices"
	"github.com/Tom5521/gotext-tools/internal/util"
)

// Entries represents a collection of Entry objects.
type Entries []Entry

func (e Entries) Equal(e2 Entries) bool {
	return util.Equal(e, e2)
}

func (e Entries) ContainsUnifiedID(uid string) bool {
	return slices.ContainsFunc(e, func(e Entry) bool { return e.UnifiedID() == uid })
}

func (e Entries) CutHeader() Entries {
	if id := e.Index("", ""); id != -1 {
		e = slices.Delete(e, id, id+1)
	}

	return e
}

func (e Entries) BestIDRatio(e1 Entry) (best, highestRatio int) {
	for i, e2 := range e {
		ratio := EntryIDMatchRatio(e1, e2)
		if ratio > highestRatio {
			best = i
			highestRatio = ratio
		}
	}
	return
}

func (e Entries) BestStrRatio(e1 Entry) (best, highestRatio int) {
	for i, e2 := range e {
		ratio := EntryStrMatchRatio(e1, e2)
		if ratio > highestRatio {
			best = i
			highestRatio = ratio
		}
	}

	return
}

func (e Entries) IndexByUnifiedID(uid string) int {
	return slices.IndexFunc(e, func(e Entry) bool {
		return e.UnifiedID() == uid
	})
}

func (e Entries) Index(id, context string) int {
	return slices.IndexFunc(e,
		func(e Entry) bool {
			return e.ID == id && e.Context == context
		},
	)
}

func (e Entries) IsSorted() bool {
	return slices.IsSortedFunc(e, CompareEntry)
}

func (e Entries) IsSortedFunc(cmp Cmp[Entry]) bool {
	return slices.IsSortedFunc(e, cmp)
}

// Sort organizes the entries by grouping them by file and sorting them by line.
func (e Entries) Sort() Entries {
	return e.SortFunc(CompareEntry)
}

func (e Entries) PrepareSorter(cmp Cmp[Entry]) func() Entries {
	return func() Entries { return e.SortFunc(cmp) }
}

func (e Entries) SortFunc(cmp Cmp[Entry]) Entries {
	slices.SortFunc(e, cmp)
	return e
}

func (e Entries) HasDuplicates() bool {
	seen := make(map[string]bool)

	for _, entry := range e {
		uid := entry.UnifiedID()
		_, seened := seen[uid]
		if seened {
			return true
		}

		seen[uid] = true
	}

	return false
}

func (e Entries) CleanObsoletes() Entries {
	return slices.DeleteFunc(e, func(e Entry) bool {
		return e.Obsolete
	})
}

// CleanDuplicates removes duplicate entries with the same ID and context, merging their locations.
func (e Entries) CleanDuplicates() Entries {
	return e.SolveFunc(func(a, b Entry) *Entry {
		a.Locations = append(a.Locations, b.Locations...)
		return &a
	})
}

// MergeFunc defines a function type that takes two Entry objects and returns a merged Entry pointer.
type MergeFunc func(a, b Entry) *Entry

// MergeUsingPriorAsBase returns a MergeFunc that prefers using the version of the entry from priorList as the base.
// If a matching UnifiedID is found in priorList, it merges that entry with 'b'.
// Otherwise, it merges 'a' and 'b' as usual using SolveMerge.
func MergeUsingPriorAsBase(priorList Entries) MergeFunc {
	return func(a, b Entry) *Entry {
		if i := priorList.IndexByUnifiedID(a.UnifiedID()); i != -1 {
			return SolveMerge(priorList[i], b)
		}
		return SolveMerge(a, b)
	}
}

// MergeAndMarkObsoleteIfNotPrioritized returns a MergeFunc that merges two entries using SolveMerge.
// If the UnifiedID of the entry is not found in the priorList, the resulting merged entry is marked as obsolete.
// This is useful for keeping non-prioritized entries while flagging them as deprecated or no longer in use.
func MergeAndMarkObsoleteIfNotPrioritized(priorList []string) MergeFunc {
	return func(a, b Entry) *Entry {
		n := SolveMerge(a, b)
		if !slices.Contains(priorList, a.UnifiedID()) {
			n.Obsolete = true
		}
		return n
	}
}

// MergeIfInPriorityList returns a MergeFunc that applies SolveMerge only if the entry's UnifiedID
// is present in the provided list of priority IDs (priorList).
// If the UnifiedID is not found in the list, the merge is skipped and nil is returned.
func MergeIfInPriorityList(priorList []string) MergeFunc {
	return func(a, b Entry) *Entry {
		if slices.Contains(priorList, a.UnifiedID()) {
			return SolveMerge(a, b)
		}
		return nil
	}
}

// MergeUsingPriorityOrFallback returns a MergeFunc that prioritizes the version from priorList if found,
// otherwise it falls back to SolveMerge.
func MergeUsingPriorityOrFallback(priorList Entries) MergeFunc {
	return func(a, b Entry) *Entry {
		if i := priorList.IndexByUnifiedID(a.UnifiedID()); i != -1 {
			return &priorList[i]
		}
		return SolveMerge(a, b)
	}
}

// SolveMerge merges two Entry objects based on certain preference criteria.
// It prefers the Entry with a higher priority according to CompareEntry.
// It combines the Locations from both entries.
// Then, it chooses the Str and Plurals fields based on CompareEntryByStr.
func SolveMerge(a, b Entry) *Entry {
	var preferred Entry

	// Choose the preferred entry based on CompareEntry.
	if CompareEntry(a, b) > 0 {
		preferred = a
	} else {
		preferred = b
	}

	// Combine the Locations from both entries.
	preferred.Locations = append(a.Locations, b.Locations...)

	// Choose the preferred Str and Plurals based on CompareEntryByStr.
	if CompareEntryByStr(a, b) > 0 {
		preferred.Str = a.Str
		preferred.Plurals = a.Plurals
	} else {
		preferred.Str = b.Str
		preferred.Plurals = b.Plurals
	}

	return &preferred
}

// SolveFunc processes a slice of Entries using a provided merging function (merger).
// It groups entries by a unified identifier (UnifiedID) and merges duplicates using the merger.
// The result is a cleaned list of Entries with duplicates resolved.
func (e Entries) SolveFunc(merger MergeFunc) Entries {
	var cleaned Entries
	seened := make(map[string]int) // Tracks indices of seen unified IDs.

	for _, entry := range e {
		uid := entry.UnifiedID()
		idIndex, ok := seened[uid]
		if ok {
			// If the ID has been seen, merge with the existing entry.
			merged := merger(cleaned[idIndex], entry)
			if merged != nil {
				cleaned[idIndex] = *merged
			}
			continue
		}

		// If not seen before, add it to the cleaned list and map.
		seened[uid] = len(cleaned)
		cleaned = append(cleaned, entry)
	}

	return cleaned
}

// Solve is a convenience method that uses the default SolveMerge function to resolve duplicates in Entries.
// It calls SolveFunc with SolveMerge as the merger function.
func (e Entries) Solve() Entries {
	return e.SolveFunc(SolveMerge)
}

func (e Entries) CleanFuzzy() Entries {
	e = slices.DeleteFunc(e, func(e Entry) bool {
		return e.IsFuzzy()
	})
	return e
}

func (e Entries) FuzzyFind(id, context string) int {
	return slices.IndexFunc(e, func(e Entry) bool {
		return util.FuzzyEqual(id, e.ID) && e.Context == context
	})
}
