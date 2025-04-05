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

func MergeWithMergeConfig(n MergeConfig) MergeOption {
	return func(mc *MergeConfig) {
		*mc = n
	}
}

func MergeWithFuzzyMatch(f bool) MergeOption {
	return func(mc *MergeConfig) {
		mc.FuzzyMatch = f
	}
}

func MergeWithKeepPreviousIDs(k bool) MergeOption {
	return func(mc *MergeConfig) {
		mc.KeepPreviousIDs = k
	}
}

func MergeWithConfig(config MergeConfig, def, ref Entries) Entries {
	def = def.Solve()

	for i, e := range def {
		if !ref.ContainsUnifiedID(e.UnifiedID()) {
			if _, ratio := ref.BestRatio(e); ratio > 50 && !e.IsFuzzy() && config.FuzzyMatch {
				e.Flags = append(e.Flags, "fuzzy")
				goto finish
			}
			if !config.KeepPreviousIDs {
				e.Obsolete = true
			}
		}

	finish:
		def[i] = e
	}
	for _, e := range ref {
		if !def.ContainsUnifiedID(e.UnifiedID()) {
			if id, ratio := def.BestRatio(e); ratio > 50 && !e.IsFuzzy() && config.FuzzyMatch {
				e.Flags = append(e.Flags, "fuzzy")

				best := def[id]
				e.Str = best.Str
				e.Plurals = best.Plurals
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
