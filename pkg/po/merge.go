package po

// SortMode defines the mode used to sort PO entries.
type SortMode int

const (
	SortByAll      SortMode = iota // Sort by all criteria.
	SortByID                       // Sort by entry ID.
	SortByFile                     // Sort by source file.
	SortByLine                     // Sort by line number.
	SortByFuzzy                    // Sort by fuzzy status.
	SortByObsolete                 // Sort by obsolete status.
)

var sortMap = map[SortMode]Cmp[Entry]{
	SortByAll:      CompareEntry,
	SortByID:       CompareEntryByID,
	SortByFile:     CompareEntryByFile,
	SortByLine:     CompareEntryByLine,
	SortByFuzzy:    CompareEntryByFuzzy,
	SortByObsolete: CompareEntryByObsolete,
}

// SortMethod returns a sorting function for the given SortMode.
func (mode SortMode) SortMethod(entries Entries) func() Entries {
	method, ok := sortMap[mode]
	if !ok {
		return entries.Sort
	}
	return entries.PrepareSorter(method)
}

// MergeConfig specifies options for merging entry sets.
type MergeConfig struct {
	FuzzyMatch      bool     // Enables fuzzy matching of entries.
	KeepPreviousIDs bool     // Retains original IDs on unmatched entries.
	Sort            bool     // Enables sorting after merge.
	SortMode        SortMode // Sort order to use if sorting is enabled.
}

// DefaultMergeConfig returns a MergeConfig with default values, optionally modified by MergeOptions.
func DefaultMergeConfig(opts ...MergeOption) MergeConfig {
	cfg := MergeConfig{
		FuzzyMatch: true,
		Sort:       true,
		SortMode:   SortByAll,
	}
	cfg.ApplyOption(opts...)
	return cfg
}

// ApplyOption applies a sequence of MergeOptions to the MergeConfig.
func (m *MergeConfig) ApplyOption(opts ...MergeOption) {
	for _, mo := range opts {
		mo(m)
	}
}

// MergeOption modifies a MergeConfig.
type MergeOption func(mc *MergeConfig)

// MergeWithSortMode returns a MergeOption that sets the SortMode.
func MergeWithSortMode(sm SortMode) MergeOption {
	return func(mc *MergeConfig) { mc.SortMode = sm }
}

// MergeWithSort returns a MergeOption that enables or disables sorting.
func MergeWithSort(s bool) MergeOption {
	return func(mc *MergeConfig) { mc.Sort = s }
}

// MergeWithMergeConfig returns a MergeOption that replaces the current config.
func MergeWithMergeConfig(n MergeConfig) MergeOption {
	return func(mc *MergeConfig) { *mc = n }
}

// MergeWithFuzzyMatch returns a MergeOption that enables or disables fuzzy matching.
func MergeWithFuzzyMatch(f bool) MergeOption {
	return func(mc *MergeConfig) { mc.FuzzyMatch = f }
}

// MergeWithKeepPreviousIDs returns a MergeOption that enables or disables retention of previous IDs.
func MergeWithKeepPreviousIDs(k bool) MergeOption {
	return func(mc *MergeConfig) { mc.KeepPreviousIDs = k }
}

// MergeWithConfig merges entries from def and ref using the provided MergeConfig.
//
// If FuzzyMatch is enabled, unmatched entries may be matched based on string similarity.
// If KeepPreviousIDs is set, original IDs are preserved in unmatched entries.
// If Sort is enabled, the resulting set is sorted using the given SortMode.
func MergeWithConfig(config MergeConfig, def, ref Entries) Entries {
	nplurals := ref.Header().Nplurals()
	def = def.Solve()

	for i, entry := range def {
		if ref.ContainsUnifiedID(entry.UnifiedID()) || entry.IsHeader() {
			continue
		}
		switch {
		case config.FuzzyMatch:
			if bestID, ratio := ref.BestIDRatio(entry); ratio >= 50 {
				entry.markAsFuzzy()
				entry.ID = ref[bestID].ID
			} else {
				entry.markAsObsolete()
			}
		case config.KeepPreviousIDs:
			entry.markAsFuzzy()
		default:
			entry.markAsObsolete()
		}
		def[i] = entry
	}

	for _, entry := range ref {
		if def.ContainsUnifiedID(entry.UnifiedID()) || entry.IsHeader() {
			continue
		}
		if config.FuzzyMatch {
			if bestID, ratio := def.BestIDRatio(entry); ratio >= 50 {
				entry.markAsFuzzy()
				best := def[bestID]
				mergeEntryStrings(&entry, best)
			}
		} else if entry.IsPlural() {
			for i := 0; i < int(nplurals); i++ {
				entry.Plurals = append(entry.Plurals, PluralEntry{i, entry.ID})
			}
		}
		def = append(def, entry)
	}

	if config.Sort {
		def = config.SortMode.SortMethod(def)()
	}

	return def
}

// mergeEntryStrings copies translation strings from source to target.
//
// If both entries are plural, plural forms are copied.
// If only one is plural, values are adapted accordingly.
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

// Merge merges def and ref using the default merge configuration,
// optionally modified by MergeOptions.
func Merge(def, ref Entries, options ...MergeOption) Entries {
	return MergeWithConfig(DefaultMergeConfig(options...), def, ref)
}
