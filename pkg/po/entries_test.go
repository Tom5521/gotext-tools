package po_test

import (
	"testing"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po"
)

func TestEntriesSort(t *testing.T) {
	tests := map[string]struct {
		input    po.Entries
		expected po.Entries
	}{
		"Empty Entries": {
			input:    po.Entries{},
			expected: po.Entries{},
		},
		"Single Entry": {
			input: po.Entries{
				{ID: "id1", Locations: []po.Location{{File: "file1.go", Line: 1}}},
			},
			expected: po.Entries{
				{ID: "id1", Locations: []po.Location{{File: "file1.go", Line: 1}}},
			},
		},
		"Multiple Entries Sorted by File and Line": {
			input: po.Entries{
				{ID: "id2", Locations: []po.Location{{File: "file2.go", Line: 5}}},
				{ID: "id1", Locations: []po.Location{{File: "file1.go", Line: 3}}},
			},
			expected: po.Entries{
				{ID: "id1", Locations: []po.Location{{File: "file1.go", Line: 3}}},
				{ID: "id2", Locations: []po.Location{{File: "file2.go", Line: 5}}},
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
		input    po.Entries
		expected po.Entries
	}{
		"Empty Entries": {
			input:    po.Entries{},
			expected: po.Entries{},
		},
		"Single Entry": {
			input: po.Entries{
				{ID: "id1", Locations: []po.Location{{File: "file1.go", Line: 1}}},
			},
			expected: po.Entries{
				{ID: "id1", Locations: []po.Location{{File: "file1.go", Line: 1}}},
			},
		},
		"Multiple Entries Sorted by File": {
			input: po.Entries{
				{ID: "id2", Locations: []po.Location{{File: "file2.go", Line: 5}}},
				{ID: "id1", Locations: []po.Location{{File: "file1.go", Line: 3}}},
			},
			expected: po.Entries{
				{ID: "id1", Locations: []po.Location{{File: "file1.go", Line: 3}}},
				{ID: "id2", Locations: []po.Location{{File: "file2.go", Line: 5}}},
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
		input    po.Entries
		expected po.Entries
	}{
		"Empty Entries": {
			input:    po.Entries{},
			expected: po.Entries{},
		},
		"Single Entry": {
			input: po.Entries{
				{ID: "id1"},
			},
			expected: po.Entries{
				{ID: "id1"},
			},
		},
		"Multiple Entries Sorted by ID": {
			input: po.Entries{
				{ID: "id2"},
				{ID: "id1"},
			},
			expected: po.Entries{
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
		input    po.Entries
		expected po.Entries
	}{
		"Empty Entries": {
			input:    po.Entries{},
			expected: po.Entries{},
		},
		"Single Entry": {
			input: po.Entries{
				{ID: "id1", Locations: []po.Location{{File: "file1.go", Line: 1}}},
			},
			expected: po.Entries{
				{ID: "id1", Locations: []po.Location{{File: "file1.go", Line: 1}}},
			},
		},
		"Multiple Entries Sorted by Line": {
			input: po.Entries{
				{ID: "id2", Locations: []po.Location{{File: "file2.go", Line: 5}}},
				{ID: "id1", Locations: []po.Location{{File: "file1.go", Line: 3}}},
			},
			expected: po.Entries{
				{ID: "id1", Locations: []po.Location{{File: "file1.go", Line: 3}}},
				{ID: "id2", Locations: []po.Location{{File: "file2.go", Line: 5}}},
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
		input    po.Entries
		expected po.Entries
	}{
		"Empty Entries": {
			input:    po.Entries{},
			expected: po.Entries{},
		},
		"No Duplicates": {
			input: po.Entries{
				{ID: "id1", Context: "ctx1"},
				{ID: "id2", Context: "ctx2"},
			},
			expected: po.Entries{
				{ID: "id1", Context: "ctx1"},
				{ID: "id2", Context: "ctx2"},
			},
		},
		"With Duplicates": {
			input: po.Entries{
				{ID: "id1", Context: "ctx1"},
				{ID: "id1", Context: "ctx1"},
			},
			expected: po.Entries{
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
		input    po.Entries
		expected po.Entries
	}{
		"Empty Entries": {
			input:    po.Entries{},
			expected: po.Entries{},
		},
		"No Conflicts": {
			input: po.Entries{
				{ID: "id1", Context: "ctx1", Str: "str1"},
				{ID: "id2", Context: "ctx2", Str: "str2"},
			},
			expected: po.Entries{
				{ID: "id1", Context: "ctx1", Str: "str1"},
				{ID: "id2", Context: "ctx2", Str: "str2"},
			},
		},
		"With Conflicts": {
			input: po.Entries{
				{ID: "id1", Context: "ctx1", Str: ""},
				{ID: "id1", Context: "ctx1", Str: "resolved"},
			},
			expected: po.Entries{
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
