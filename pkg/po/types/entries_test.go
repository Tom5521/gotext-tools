package types_test

import (
	"testing"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po/types"
)

func TestEntriesSort(t *testing.T) {
	tests := map[string]struct {
		input    types.Entries
		expected types.Entries
	}{
		"Empty Entries": {
			input:    types.Entries{},
			expected: types.Entries{},
		},
		"Single Entry": {
			input: types.Entries{
				{ID: "id1", Locations: []types.Location{{File: "file1.go", Line: 1}}},
			},
			expected: types.Entries{
				{ID: "id1", Locations: []types.Location{{File: "file1.go", Line: 1}}},
			},
		},
		"Multiple Entries Sorted by File and Line": {
			input: types.Entries{
				{ID: "id2", Locations: []types.Location{{File: "file2.go", Line: 5}}},
				{ID: "id1", Locations: []types.Location{{File: "file1.go", Line: 3}}},
			},
			expected: types.Entries{
				{ID: "id1", Locations: []types.Location{{File: "file1.go", Line: 3}}},
				{ID: "id2", Locations: []types.Location{{File: "file2.go", Line: 5}}},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.input.Sort()
			if !util.Equal(result, tc.expected) {
				t.Errorf("Test failed: %s\nExpected: %+v\nGot: %+v", name, tc.expected, result)
			}
		})
	}
}

func TestEntriesSortByFile(t *testing.T) {
	tests := map[string]struct {
		input    types.Entries
		expected types.Entries
	}{
		"Empty Entries": {
			input:    types.Entries{},
			expected: types.Entries{},
		},
		"Single Entry": {
			input: types.Entries{
				{ID: "id1", Locations: []types.Location{{File: "file1.go", Line: 1}}},
			},
			expected: types.Entries{
				{ID: "id1", Locations: []types.Location{{File: "file1.go", Line: 1}}},
			},
		},
		"Multiple Entries Sorted by File": {
			input: types.Entries{
				{ID: "id2", Locations: []types.Location{{File: "file2.go", Line: 5}}},
				{ID: "id1", Locations: []types.Location{{File: "file1.go", Line: 3}}},
			},
			expected: types.Entries{
				{ID: "id1", Locations: []types.Location{{File: "file1.go", Line: 3}}},
				{ID: "id2", Locations: []types.Location{{File: "file2.go", Line: 5}}},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.input.SortByFile()
			if !util.Equal(result, tc.expected) {
				t.Errorf("Test failed: %s\nExpected: %+v\nGot: %+v", name, tc.expected, result)
			}
		})
	}
}

func TestEntriesSortByID(t *testing.T) {
	tests := map[string]struct {
		input    types.Entries
		expected types.Entries
	}{
		"Empty Entries": {
			input:    types.Entries{},
			expected: types.Entries{},
		},
		"Single Entry": {
			input: types.Entries{
				{ID: "id1"},
			},
			expected: types.Entries{
				{ID: "id1"},
			},
		},
		"Multiple Entries Sorted by ID": {
			input: types.Entries{
				{ID: "id2"},
				{ID: "id1"},
			},
			expected: types.Entries{
				{ID: "id1"},
				{ID: "id2"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.input.SortByID()
			if !util.Equal(result, tc.expected) {
				t.Errorf("Test failed: %s\nExpected: %+v\nGot: %+v", name, tc.expected, result)
			}
		})
	}
}

func TestEntriesSortByLine(t *testing.T) {
	tests := map[string]struct {
		input    types.Entries
		expected types.Entries
	}{
		"Empty Entries": {
			input:    types.Entries{},
			expected: types.Entries{},
		},
		"Single Entry": {
			input: types.Entries{
				{ID: "id1", Locations: []types.Location{{File: "file1.go", Line: 1}}},
			},
			expected: types.Entries{
				{ID: "id1", Locations: []types.Location{{File: "file1.go", Line: 1}}},
			},
		},
		"Multiple Entries Sorted by Line": {
			input: types.Entries{
				{ID: "id2", Locations: []types.Location{{File: "file2.go", Line: 5}}},
				{ID: "id1", Locations: []types.Location{{File: "file1.go", Line: 3}}},
			},
			expected: types.Entries{
				{ID: "id1", Locations: []types.Location{{File: "file1.go", Line: 3}}},
				{ID: "id2", Locations: []types.Location{{File: "file2.go", Line: 5}}},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.input.SortByLine()
			if !util.Equal(result, tc.expected) {
				t.Errorf("Test failed: %s\nExpected: %+v\nGot: %+v", name, tc.expected, result)
			}
		})
	}
}

func TestEntriesCleanDuplicates(t *testing.T) {
	tests := map[string]struct {
		input    types.Entries
		expected types.Entries
	}{
		"Empty Entries": {
			input:    types.Entries{},
			expected: types.Entries{},
		},
		"No Duplicates": {
			input: types.Entries{
				{ID: "id1", Context: "ctx1"},
				{ID: "id2", Context: "ctx2"},
			},
			expected: types.Entries{
				{ID: "id1", Context: "ctx1"},
				{ID: "id2", Context: "ctx2"},
			},
		},
		"With Duplicates": {
			input: types.Entries{
				{ID: "id1", Context: "ctx1"},
				{ID: "id1", Context: "ctx1"},
			},
			expected: types.Entries{
				{ID: "id1", Context: "ctx1"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.input.CleanDuplicates()
			if !util.Equal(result, tc.expected) {
				t.Errorf("Test failed: %s\nExpected: %+v\nGot: %+v", name, tc.expected, result)
			}
		})
	}
}

func TestEntriesSolve(t *testing.T) {
	tests := map[string]struct {
		input    types.Entries
		expected types.Entries
	}{
		"Empty Entries": {
			input:    types.Entries{},
			expected: types.Entries{},
		},
		"No Conflicts": {
			input: types.Entries{
				{ID: "id1", Context: "ctx1", Str: "str1"},
				{ID: "id2", Context: "ctx2", Str: "str2"},
			},
			expected: types.Entries{
				{ID: "id1", Context: "ctx1", Str: "str1"},
				{ID: "id2", Context: "ctx2", Str: "str2"},
			},
		},
		"With Conflicts": {
			input: types.Entries{
				{ID: "id1", Context: "ctx1", Str: ""},
				{ID: "id1", Context: "ctx1", Str: "resolved"},
			},
			expected: types.Entries{
				{ID: "id1", Context: "ctx1", Str: "resolved"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.input.Solve()
			if !util.Equal(result, tc.expected) {
				t.Errorf("Test failed: %s\nExpected: %+v\nGot: %+v", name, tc.expected, result)
			}
		})
	}
}
