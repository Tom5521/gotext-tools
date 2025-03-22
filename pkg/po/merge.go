package po

import "strings"

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

func (f File) MergeWithConfig(config MergeConfig, files ...*File) *File {
	names := []string{f.Name}
	for _, file := range files {
		names = append(names, file.Name)
		f.Entries = append(f.Entries, file.Entries...)
	}
	f.Name = strings.Join(names, "_")

	if config.FuzzyMatch {
		f.Entries = f.Entries.FuzzySolve()
	} else {
		f.Entries = f.Entries.Solve()
	}
	if config.Sort {
		f.Entries = config.SortMode.SortMethod(f.Entries)()
	}

	return &f
}

func (f File) MergeWithOptions(files []*File, options ...MergeOption) *File {
	var cfg MergeConfig
	cfg.ApplyOption(options...)
	return f.MergeWithConfig(cfg, files...)
}

func (f File) Merge(files ...*File) *File {
	return f.MergeWithConfig(DefaultMergeConfig(), files...)
}
