package po

import "strings"

type SortMode int

const (
	SortByAll SortMode = iota
	SortByID
	SortByFile
	SortByLine
	SortByFuzzy
)

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
		switch config.SortMode {
		case SortByID:
			f.Entries = f.Entries.SortByID()
		case SortByFile:
			f.Entries = f.Entries.SortByFile()
		case SortByLine:
			f.Entries = f.Entries.SortByLine()
		case SortByFuzzy:
			f.Entries = f.Entries.SortByFuzzy()
		case SortByAll:
			fallthrough
		default:
			f.Entries = f.Entries.Sort()
		}
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
