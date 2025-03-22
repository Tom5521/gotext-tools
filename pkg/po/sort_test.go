package po_test

import (
	"math/rand"
	"slices"
	"testing"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po"
)

func TestSort(t *testing.T) {
	tests := []struct {
		name     string
		cmp      func(a, b po.Entry) int
		expected po.Entries
	}{
		{
			"Sort",
			po.CompareEntry,
			po.Entries{
				{Locations: []po.Location{{File: "a", Line: 1}}},
				{Locations: []po.Location{{File: "a", Line: 2}}},
				{Locations: []po.Location{{File: "b", Line: 10}}},
				{Locations: []po.Location{{File: "c", Line: 200}}},
				{Obsolete: true},
				{Flags: []string{"fuzzy"}},
			},
		},
		{
			"SortByLine",
			po.CompareEntryByLine,
			po.Entries{
				{Locations: []po.Location{{Line: 1}}},
				{Locations: []po.Location{{Line: 2}}},
				{Locations: []po.Location{{Line: 3}}},
				{Locations: []po.Location{{Line: 4}}},
				{Locations: []po.Location{{Line: 5}}},
			},
		},
		{
			"SortByFile",
			po.CompareEntryByFile,
			po.Entries{
				{Locations: []po.Location{{File: "a.txt"}}},
				{Locations: []po.Location{{File: "b.txt"}}},
				{Locations: []po.Location{{File: "c.txt"}}},
				{Locations: []po.Location{{File: "d.txt"}}},
				{Locations: []po.Location{{File: "e.txt"}}},
				{Locations: []po.Location{{File: "f.txt"}}},
			},
		},
		{
			"SortByID",
			po.CompareEntryByID,
			po.Entries{
				{ID: "a"},
				{ID: "b"},
				{ID: "c"},
				{ID: "d"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sorted := slices.Clone(test.expected)
			rand.Shuffle(len(sorted), func(i, j int) {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			})

			slices.SortFunc(sorted, test.cmp)

			if !util.Equal(sorted, test.expected) {
				t.Fail()
			}
		})
	}
}

func BenchmarkSortEntries(b *testing.B) {
	entries := po.Entries{
		{
			ID:      "Apple",
			Context: "USA",
			Plural:  "Apples",
			Plurals: po.PluralEntries{
				{0, "Manzana"},
				{1, "Manzanas"},
			},
		},
		{ID: "Hi", Str: "Hola", Context: "casual"},
		{ID: "", Str: ""},
		{ID: "How are you?", Str: "Como estás?"},
	}

	rand.Shuffle(len(entries), func(i, j int) {
		entries[i], entries[j] = entries[j], entries[i]
	})

	tests := []struct {
		name   string
		method func() po.Entries
	}{
		{"Sort", entries.Sort},
		{"SortByFile", entries.SortByFile},
		{"SortByID", entries.SortByID},
		{"SortByLine", entries.SortByLine},
		{"SortByFuzzy", entries.SortByFuzzy},
	}

	for _, t := range tests {
		b.Run(t.name, func(b *testing.B) {
			t.method()
		})
	}
}

func BenchmarkEntriesSolve(b *testing.B) {
	entries := po.Entries{
		{
			ID:      "Apple",
			Context: "USA",
			Plural:  "Apples",
			Plurals: po.PluralEntries{
				{0, "Manzana"},
				{1, "Manzanas"},
			},
		},
		{ID: "Hi", Str: "Hola", Context: "casual"},
		{ID: "", Str: ""},
		{ID: "How are you?", Str: "Como estás?"},
	}

	rand.Shuffle(len(entries), func(i, j int) {
		entries[i], entries[j] = entries[j], entries[i]
	})

	tests := []struct {
		name   string
		method func() po.Entries
	}{
		{"Solve", entries.Solve},
		{"FuzzySolve", entries.FuzzySolve},
	}

	for _, t := range tests {
		b.Run(t.name, func(b *testing.B) {
			t.method()
		})
	}
}
