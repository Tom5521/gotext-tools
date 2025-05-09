package po_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/Tom5521/gotext-tools/v2/internal/slices"
	"github.com/Tom5521/gotext-tools/v2/internal/util"
	"github.com/Tom5521/gotext-tools/v2/pkg/po"
	"github.com/Tom5521/gotext-tools/v2/pkg/po/compile"
)

func TestSolve(t *testing.T) {
	input := po.Entries{
		{ID: "Hello", Str: "World"},
		{ID: "Hello2", Str: "World2"},
		{ID: "Hello3", Str: "World3"},
		{ID: "Hello"},
		{ID: "Hello2"},
		{ID: "Hello3"},
	}
	expected := po.Entries{
		{ID: "Hello", Str: "World"},
		{ID: "Hello2", Str: "World2"},
		{ID: "Hello3", Str: "World3"},
	}

	solved := input.Solve()

	if !util.Equal(solved, expected) {
		fmt.Println(
			compile.PoToString(solved, compile.PoWithOmitHeader(true)),
		)
		fmt.Println(util.NamedDiff("solved", "expected", solved, expected))

		t.Fail()
	}
}

func TestSort(t *testing.T) {
	tests := []struct {
		name     string
		cmp      po.Cmp[po.Entry]
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

			sorted = sorted.SortFunc(test.cmp)

			if !util.Equal(sorted, test.expected) {
				fmt.Println(util.NamedDiff("sorted", "expected", sorted, test.expected))
				t.Fail()
			}
		})
	}
}
