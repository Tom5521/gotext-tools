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

func (mode SortMode) SortMethod(entries Entries) func() Entries {
	method, ok := map[SortMode]func() Entries{
		SortByAll:      entries.Sort,
		SortByID:       entries.PrepareSorter(CompareEntryByID),
		SortByFile:     entries.PrepareSorter(CompareEntryByFile),
		SortByLine:     entries.PrepareSorter(CompareEntryByLine),
		SortByFuzzy:    entries.PrepareSorter(CompareEntryByFuzzy),
		SortByObsolete: entries.PrepareSorter(CompareEntryByObsolete),
	}[mode]

	if !ok {
		return entries.Sort
	}

	return method
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

func MergeWithConfig(config MergeConfig, def, ref Entries) Entries {
	def = def.Solve()

	for i, e := range def {
		if !ref.ContainsUnifiedID(e.UnifiedID()) && !e.IsHeader() {
			if config.FuzzyMatch {
				if bestID, ratio := ref.BestIDRatio(e); ratio >= 50 {
					e.markAsFuzzy()
					best := ref[bestID]
					e.ID = best.ID
				} else {
					e.markAsObsolete()
				}
			} else {
				if config.KeepPreviousIDs {
					e.markAsFuzzy()
				} else {
					e.markAsObsolete()
				}
			}

			def[i] = e
		}
	}

	for _, e := range ref {
		if !def.ContainsUnifiedID(e.UnifiedID()) && !e.IsHeader() {
			if config.FuzzyMatch {
				if id, ratio := def.BestIDRatio(e); ratio >= 50 {
					e.markAsFuzzy()
					best := def[id]
					switch {
					case e.IsPlural() && best.IsPlural():
						e.Plurals = best.Plurals
					case e.IsPlural() && !best.IsPlural():
						if len(e.Plurals) == 0 {
							e.Plurals = append(e.Plurals, PluralEntry{})
						}
						e.Plurals[0].Str = best.Str
					case !e.IsPlural() && best.IsPlural():
						if len(best.Plurals) > 0 {
							e.Str = best.Plurals[0].Str
						}
					case !e.IsPlural() && !best.IsPlural():
						e.Str = best.Str
					}
				}
			} else if e.IsPlural() {
				e.Plurals = PluralEntries{
					{0, e.ID},
					{1, e.ID},
				}
			}

			def = append(def, e)
		}
	}

	if config.Sort {
		def = config.SortMode.SortMethod(def)()
	}

	return def
}

func Merge(def, ref Entries, options ...MergeOption) Entries {
	return MergeWithConfig(DefaultMergeConfig(options...), def, ref)
}
