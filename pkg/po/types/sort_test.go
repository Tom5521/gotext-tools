package types_test

import (
	"testing"

	"github.com/Tom5521/xgotext/pkg/po/types"
)

func TestSortEntriesByID(t *testing.T) {
	entries := []types.Entry{
		{ID: "b"},
		{ID: "a"},
		{ID: "c"},
	}

	expected := []types.Entry{
		{ID: "a"},
		{ID: "b"},
		{ID: "c"},
	}

	types.SortEntriesByID(entries)

	for i, entry := range entries {
		if entry.ID != expected[i].ID {
			t.Errorf("Expected ID %s at index %d, got %s", expected[i].ID, i, entry.ID)
		}
	}
}

func TestSortEntriesByLine(t *testing.T) {
	entries := []types.Entry{
		{Locations: []types.Location{{Line: 3}}},
		{Locations: []types.Location{{Line: 1}}},
		{Locations: []types.Location{{Line: 2}}},
	}

	expected := []types.Entry{
		{Locations: []types.Location{{Line: 1}}},
		{Locations: []types.Location{{Line: 2}}},
		{Locations: []types.Location{{Line: 3}}},
	}

	types.SortEntriesByLine(entries)

	for i, entry := range entries {
		if entry.Locations[0].Line != expected[i].Locations[0].Line {
			t.Errorf(
				"Expected line %d at index %d, got %d",
				expected[i].Locations[0].Line,
				i,
				entry.Locations[0].Line,
			)
		}
	}
}

func TestSortEntriesByFile(t *testing.T) {
	entries := []types.Entry{
		{Locations: []types.Location{{File: "b.txt"}}},
		{Locations: []types.Location{{File: "a.txt"}}},
		{Locations: []types.Location{{File: "c.txt"}}},
	}

	expected := []types.Entry{
		{Locations: []types.Location{{File: "a.txt"}}},
		{Locations: []types.Location{{File: "b.txt"}}},
		{Locations: []types.Location{{File: "c.txt"}}},
	}

	types.SortEntriesByFile(entries)

	for i, entry := range entries {
		if entry.Locations[0].File != expected[i].Locations[0].File {
			t.Errorf(
				"Expected file %s at index %d, got %s",
				expected[i].Locations[0].File,
				i,
				entry.Locations[0].File,
			)
		}
	}
}

func TestSortEntriesByFileAndLine(t *testing.T) {
	entries := []types.Entry{
		{Locations: []types.Location{{File: "b.txt", Line: 2}}},
		{Locations: []types.Location{{File: "a.txt", Line: 2}}},
		{Locations: []types.Location{{File: "a.txt", Line: 1}}},
		{Locations: []types.Location{{File: "c.txt", Line: 1}}},
	}

	expected := []types.Entry{
		{Locations: []types.Location{{File: "a.txt", Line: 1}}},
		{Locations: []types.Location{{File: "a.txt", Line: 2}}},
		{Locations: []types.Location{{File: "b.txt", Line: 2}}},
		{Locations: []types.Location{{File: "c.txt", Line: 1}}},
	}

	types.SortEntriesByFileAndLine(entries)

	for i, entry := range entries {
		if entry.Locations[0].File != expected[i].Locations[0].File ||
			entry.Locations[0].Line != expected[i].Locations[0].Line {
			t.Errorf("Expected file %s and line %d at index %d, got file %s and line %d",
				expected[i].Locations[0].File, expected[i].Locations[0].Line, i,
				entry.Locations[0].File, entry.Locations[0].Line)
		}
	}
}
