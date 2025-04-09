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
		SortByID:       entries.SortByID,
		SortByFile:     entries.SortByFile,
		SortByLine:     entries.SortByLine,
		SortByFuzzy:    entries.SortByFuzzy,
		SortByObsolete: entries.SortByObsolete,
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
	// TODO: Finish this.
	// Compendium      Entries
}

func DefaultMergeConfig() MergeConfig {
	return MergeConfig{
		FuzzyMatch: true,
		Sort:       true,
		SortMode:   SortByAll,
	}
}

func (m *MergeConfig) ApplyOption(opts ...MergeOption) {
	for _, mo := range opts {
		mo(m)
	}
}

type MergeOption func(mc *MergeConfig)

// TODO: Finish this.
// func MergeWithCompendium(compendium Entries) MergeOption {
// 	return func(mc *MergeConfig) {
// 		mc.Compendium = compendium
// 	}
// }

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
		if !ref.ContainsUnifiedID(e.UnifiedID()) {
			if config.FuzzyMatch {
				if _, ratio := ref.BestRatio(e); ratio >= 20 {
					e.markAsFuzzy()
				} else {
					e.markAsObsolete()
				}
			} else {
				if !config.KeepPreviousIDs {
					e.markAsObsolete()
				} else {
					e.markAsFuzzy()
				}
			}

			def[i] = e
		}
	}

	for _, e := range ref {
		if !def.ContainsUnifiedID(e.UnifiedID()) {
			e.markAsFuzzy()

			if config.FuzzyMatch {
				if id, ratio := def.BestRatio(e); ratio >= 20 {
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
	cfg := DefaultMergeConfig()
	cfg.ApplyOption(options...)
	return MergeWithConfig(cfg, def, ref)
}
