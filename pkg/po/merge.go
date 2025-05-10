package po

type SortMode int

const (
	SortByAll SortMode = iota
	SortByID
	SortByFile
	SortByLine
	SortByFuzzy
	SortByObsolete
)

var sortMap = map[SortMode]Cmp[Entry]{
	SortByAll:      CompareEntry,
	SortByID:       CompareEntryByID,
	SortByFile:     CompareEntryByFile,
	SortByLine:     CompareEntryByLine,
	SortByFuzzy:    CompareEntryByFuzzy,
	SortByObsolete: CompareEntryByObsolete,
}

func (mode SortMode) SortMethod(entries Entries) func() Entries {
	method, ok := sortMap[mode]

	if !ok {
		return entries.Sort
	}

	return entries.PrepareSorter(method)
}

type MergeConfig struct {
	FuzzyMatch      bool
	KeepPreviousIDs bool
	Sort            bool
	SortMode        SortMode
}

func DefaultMergeConfig(opts ...MergeOption) MergeConfig {
	cfg := MergeConfig{
		FuzzyMatch: true,
		Sort:       true,
		SortMode:   SortByAll,
	}
	cfg.ApplyOption(opts...)

	return cfg
}

func (m *MergeConfig) ApplyOption(opts ...MergeOption) {
	for _, mo := range opts {
		mo(m)
	}
}

type MergeOption func(mc *MergeConfig)

func MergeWithSortMode(sm SortMode) MergeOption {
	return func(mc *MergeConfig) { mc.SortMode = sm }
}

func MergeWithSort(s bool) MergeOption {
	return func(mc *MergeConfig) { mc.Sort = s }
}

func MergeWithMergeConfig(n MergeConfig) MergeOption {
	return func(mc *MergeConfig) { *mc = n }
}

func MergeWithFuzzyMatch(f bool) MergeOption {
	return func(mc *MergeConfig) { mc.FuzzyMatch = f }
}

func MergeWithKeepPreviousIDs(k bool) MergeOption {
	return func(mc *MergeConfig) { mc.KeepPreviousIDs = k }
}

// MergeWithConfig merges two sets of entries (`def` and `ref`) according to the provided `config`.
//
// It performs the following steps:
// 1. Solves the `def` entries before processing.
// 2. Iterates over `def` to mark entries as fuzzy or obsolete if they are not present in `ref`.
// 3. Iterates over `ref` to add missing entries into `def`, optionally applying fuzzy match logic.
// 4. Optionally sorts the result using the provided sort configuration.
//
// Parameters:
// - config: MergeConfig object that controls fuzzy matching, ID retention, and sorting behavior.
// - def:    The default entry list to be merged and returned.
// - ref:    The reference entry list to be compared against.
//
// Returns:
// - Entries: A merged list of entries according to the merge configuration.
func MergeWithConfig(config MergeConfig, def, ref Entries) Entries {
	nplurals := ref.Header().Nplurals()
	def = def.Solve()

	// Step 1: Process entries from `def` that are not present in `ref`
	for i, e := range def {
		if ref.ContainsUnifiedID(e.UnifiedID()) || e.IsHeader() {
			continue
		}

		switch {
		case config.FuzzyMatch:
			// Attempt fuzzy match
			if bestID, ratio := ref.BestIDRatio(e); ratio >= 50 {
				e.markAsFuzzy()
				e.ID = ref[bestID].ID
			} else {
				e.markAsObsolete()
			}
		case config.KeepPreviousIDs:
			// Retain the original ID if allowed
			e.markAsFuzzy()
		default:
			e.markAsObsolete()
		}

		def[i] = e
	}

	// Step 2: Add entries from `ref` that are missing in `def`
	for _, e := range ref {
		if def.ContainsUnifiedID(e.UnifiedID()) || e.IsHeader() {
			continue
		}

		if config.FuzzyMatch {
			// Attempt fuzzy match and merge content
			if bestID, ratio := def.BestIDRatio(e); ratio >= 50 {
				e.markAsFuzzy()
				best := def[bestID]
				mergeEntryStrings(&e, best)
			}
		} else if e.IsPlural() {
			// If no fuzzy match, create default plural entries
			for i := 0; i < int(nplurals); i++ {
				e.Plurals = append(e.Plurals, PluralEntry{i, e.ID})
			}
		}

		def = append(def, e)
	}

	// Step 3: Optionally sort the final result
	if config.Sort {
		def = config.SortMode.SortMethod(def)()
	}

	return def
}

// mergeEntryStrings merges the string values of two entries, adapting plural/singular forms appropriately.
//
// If both entries are plural, copies all plural forms.
// If one is singular and the other plural, tries to adapt the string accordingly.
//
// Parameters:
// - target: The entry that will be modified.
// - source: The entry whose values will be used as the source.
func mergeEntryStrings(target *Entry, source Entry) {
	switch {
	case target.IsPlural() && source.IsPlural():
		target.Plurals = source.Plurals
	case target.IsPlural() && !source.IsPlural():
		if len(target.Plurals) == 0 {
			target.Plurals = make(PluralEntries, 1)
		}
		target.Plurals[0].Str = source.Str
	case !target.IsPlural() && source.IsPlural():
		if len(source.Plurals) > 0 {
			target.Str = source.Plurals[0].Str
		}
	case !target.IsPlural() && !source.IsPlural():
		target.Str = source.Str
	}
}

func Merge(def, ref Entries, options ...MergeOption) Entries {
	return MergeWithConfig(DefaultMergeConfig(options...), def, ref)
}
