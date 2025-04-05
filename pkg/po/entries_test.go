package po_test

import (
	"fmt"
	"math/rand"
	"slices"
	"testing"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po"
	"github.com/Tom5521/xgotext/pkg/po/compiler"
	"github.com/kr/pretty"
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
			compiler.NewPo(&po.File{Entries: solved}, compiler.PoWithOmitHeader(true)).ToString(),
		)
		for _, d := range pretty.Diff(solved, expected) {
			fmt.Println(d)
		}
		t.Fail()
	}
}

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
