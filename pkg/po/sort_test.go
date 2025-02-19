package po_test

import (
	"testing"

	"github.com/Tom5521/xgotext/pkg/po"
)

func TestSortEntriesByID(t *testing.T) {
	entries := po.Entries{
		{ID: "b"},
		{ID: "a"},
		{ID: "c"},
	}

	expected := po.Entries{
		{ID: "a"},
		{ID: "b"},
		{ID: "c"},
	}

	entries = entries.SortByID()

	for i, entry := range entries {
		if entry.ID != expected[i].ID {
			t.Errorf("Expected ID %s at index %d, got %s", expected[i].ID, i, entry.ID)
		}
	}
}

func TestSortEntriesByLine(t *testing.T) {
	entries := po.Entries{
		{Locations: []po.Location{{Line: 3}}},
		{Locations: []po.Location{{Line: 1}}},
		{Locations: []po.Location{{Line: 2}}},
	}

	expected := po.Entries{
		{Locations: []po.Location{{Line: 1}}},
		{Locations: []po.Location{{Line: 2}}},
		{Locations: []po.Location{{Line: 3}}},
	}

	entries = entries.SortByLine()

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
	entries := po.Entries{
		{Locations: []po.Location{{File: "b.txt"}}},
		{Locations: []po.Location{{File: "a.txt"}}},
		{Locations: []po.Location{{File: "c.txt"}}},
	}

	expected := po.Entries{
		{Locations: []po.Location{{File: "a.txt"}}},
		{Locations: []po.Location{{File: "b.txt"}}},
		{Locations: []po.Location{{File: "c.txt"}}},
	}

	entries = entries.SortByFile()

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

func TestSortEntries(t *testing.T) {
	entries := po.Entries{
		{Locations: []po.Location{{File: "b.txt", Line: 2}}},
		{Locations: []po.Location{{File: "a.txt", Line: 2}}},
		{Locations: []po.Location{{File: "a.txt", Line: 1}}},
		{Locations: []po.Location{{File: "c.txt", Line: 1}}},
	}

	expected := po.Entries{
		{Locations: []po.Location{{File: "a.txt", Line: 1}}},
		{Locations: []po.Location{{File: "a.txt", Line: 2}}},
		{Locations: []po.Location{{File: "b.txt", Line: 2}}},
		{Locations: []po.Location{{File: "c.txt", Line: 1}}},
	}

	entries = entries.Sort()

	for i, entry := range entries {
		if entry.Locations[0].File != expected[i].Locations[0].File ||
			entry.Locations[0].Line != expected[i].Locations[0].Line {
			t.Errorf("Expected file %s and line %d at index %d, got file %s and line %d",
				expected[i].Locations[0].File, expected[i].Locations[0].Line, i,
				entry.Locations[0].File, entry.Locations[0].Line)
		}
	}
}
