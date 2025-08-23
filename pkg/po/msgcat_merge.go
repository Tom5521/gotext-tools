package po

import (
	"fmt"
	"math"
	"strings"

	"github.com/Tom5521/gotext-tools/v2/internal/slices"
)

const msgcatFileTemplate = "#-#-#-#-#  %s  #-#-#-#-#"

func DefaultMsgcatMergeConfig(opts ...MsgcatMergeOption) MsgcatMergeConfig {
	cfg := MsgcatMergeConfig{
		LessThan: math.MaxInt, // Infinite.
	}

	cfg.ApplyOption(opts...)

	return cfg
}

type MsgcatMergeConfig struct {
	LessThan uint
	MoreThan uint
	UseFirst bool
}

func (mmc *MsgcatMergeConfig) ApplyOption(opts ...MsgcatMergeOption) {
	for _, mmo := range opts {
		mmo(mmc)
	}
}

type MsgcatMergeOption func(mc *MsgcatMergeConfig)

func MsgcatMergeWithConfig(mmc MsgcatMergeConfig) MsgcatMergeOption {
	return func(mc *MsgcatMergeConfig) { *mc = mmc }
}

func MsgcatMergeWithLessThan(l uint) MsgcatMergeOption {
	return func(mc *MsgcatMergeConfig) { mc.LessThan = l }
}

func MsgcatMergeWithMoreThan(m uint) MsgcatMergeOption {
	return func(mc *MsgcatMergeConfig) { mc.MoreThan = m }
}

func MsgcatMergeWithUseFirst(uf bool) MsgcatMergeOption {
	return func(mc *MsgcatMergeConfig) { mc.UseFirst = uf }
}

func MsgcatMergeFiles(files []*File, opts ...MsgcatMergeOption) Entries {
	process := newMsgcatMergeProcess(opts...)

	for _, file := range files {
		for _, entry := range file.Entries {
			process.handleEntry(entry, file.Name)
		}
	}

	return slices.DeleteFunc(process.entriesSlice, func(entry Entry) bool {
		info := process.entriesMap[entry.UnifiedID()]

		if info.timesFound >= process.config.LessThan {
			return true
		}
		if info.timesFound <= process.config.MoreThan {
			return true
		}

		return false
	})
}

func newMsgcatMergeProcess(opts ...MsgcatMergeOption) *msgcatMergeProcess {
	return &msgcatMergeProcess{
		map[string]msgcatEntryInfo{},
		Entries{},
		DefaultMsgcatMergeConfig(opts...),
	}
}

type msgcatMergeProcess struct {
	entriesMap   map[string]msgcatEntryInfo
	entriesSlice Entries
	config       MsgcatMergeConfig
}

// msgcatEntryInfo tracks metadata about merged entries.
type msgcatEntryInfo struct {
	index      int    // Index in the entries slice
	file       string // First file where the entry was found
	timesFound uint
}

// handleEntry handles an individual PO entry including merging logic.
func (m *msgcatMergeProcess) handleEntry(entry Entry, filename string) {
	uid := entry.UnifiedID()
	original, duplicated := m.entriesMap[uid]

	if duplicated {
		original.timesFound++
		m.entriesMap[uid] = original

		if m.config.UseFirst {
			return
		}
		m.entriesSlice[original.index] = msgcatMergeEntry(
			m.entriesSlice[original.index],
			entry,
			original.file,
			filename,
		)
		return
	}

	m.entriesMap[uid] = msgcatEntryInfo{
		index:      len(m.entriesSlice),
		file:       filename,
		timesFound: 1,
	}
	m.entriesSlice = append(m.entriesSlice, entry)
}

// msgcatMergeEntry merges two conflicting PO entries.
func msgcatMergeEntry(
	originalEntry Entry,
	newEntry Entry,
	originalFilename, newFilename string,
) Entry {
	mergedEntry := Entry{
		Flags: msgcatMergeComments(
			originalEntry.Flags,
			newEntry.Flags,
			originalFilename,
			newFilename,
		),
		Comments: msgcatMergeComments(
			originalEntry.Comments,
			newEntry.Comments,
			originalFilename,
			newFilename,
		),
		ExtractedComments: msgcatMergeComments(
			originalEntry.ExtractedComments,
			newEntry.ExtractedComments,
			originalFilename,
			newFilename,
		),
		Previous: msgcatMergeComments(
			originalEntry.Previous,
			newEntry.Previous,
			originalFilename,
			newFilename,
		),

		ID:      originalEntry.ID,
		Context: originalEntry.Context,
		Plural:  originalEntry.Plural,
		Plurals: msgcatMergePlurals(
			originalEntry,
			newEntry,
			originalFilename,
			newFilename,
		),
		Str: msgcatMergeStrings(
			originalEntry.Str,
			newEntry.Str,
			originalFilename,
			newFilename,
		),
		Locations: append(
			originalEntry.Locations,
			newEntry.Locations...,
		),
	}

	if !mergedEntry.IsFuzzy() &&
		(originalEntry.FullUnifiedID() != newEntry.FullUnifiedID() ||
			originalEntry.UnifiedStr() != newEntry.UnifiedStr()) {
		mergedEntry.Flags = append(mergedEntry.Flags, "fuzzy")
	}

	return mergedEntry
}

// msgcatMergePlurals handles merging of plural form strings.
func msgcatMergePlurals(
	originalEntry, newEntry Entry,
	originalFilename, newFilename string,
) PluralEntries {
	if !originalEntry.IsPlural() {
		return originalEntry.Plurals
	}

	mergedEntries := make(PluralEntries, 0, len(originalEntry.Plurals)+len(newEntry.Plurals))

	for i, original := range originalEntry.Plurals {
		if i >= len(newEntry.Plurals) {
			mergedEntries = append(mergedEntries, original)
			continue
		}

		newPlural := newEntry.Plurals[i]
		mergedEntries = append(mergedEntries, PluralEntry{
			ID: original.ID,
			Str: msgcatMergeStrings(
				original.Str, newPlural.Str,
				originalFilename, newFilename,
			),
		})
	}

	return mergedEntries
}

func msgcatMergeStrings(originalStr, newStr string, originalFilename, newFilename string) string {
	if originalStr == newStr {
		return originalStr
	}
	if originalStr == "" {
		return newStr
	}
	if newStr == "" {
		return originalStr
	}

	var builder strings.Builder
	fmt.Fprintf(&builder, msgcatFileTemplate+"\n", originalFilename)
	fmt.Fprintln(&builder, originalStr)
	fmt.Fprintf(&builder, msgcatFileTemplate+"\n", newFilename)
	fmt.Fprint(&builder, newStr)
	return builder.String()
}

func msgcatMergeComments(
	originalComments, newComments []string,
	originalFilename, newFilename string,
) []string {
	if slices.Equal(originalComments, newComments) {
		return originalComments
	}

	if len(originalComments) == 0 {
		return append(
			[]string{fmt.Sprintf(msgcatFileTemplate, newFilename)}, newComments...,
		)
	}
	if len(newComments) == 0 {
		return append(
			[]string{fmt.Sprintf(msgcatFileTemplate, originalFilename)},
			originalComments...,
		)
	}

	return append(
		append(
			[]string{fmt.Sprintf(msgcatFileTemplate, originalFilename)},
			originalComments...,
		),
		append(
			[]string{fmt.Sprintf(msgcatFileTemplate, newFilename)},
			newComments...,
		)...,
	)
}
