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
				{Flags: []string{"fuzzy"}},
				{Obsolete: true},
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
