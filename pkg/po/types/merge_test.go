package types

import (
	"testing"
)

func TestMergeFiles(t *testing.T) {
	file1 := &File{
		Name:     "file1",
		Header:   Header{},
		Nplurals: 2,
		Entries: Entries{
			{ID: "id1", Str: "str1", Locations: []Location{{File: "file1.go", Line: 10}}},
			{ID: "id2", Str: "str2", Locations: []Location{{File: "file1.go", Line: 20}}},
		},
	}

	file2 := &File{
		Name:     "file2",
		Header:   Header{},
		Nplurals: 3,
		Entries: Entries{
			{ID: "id3", Str: "str3", Locations: []Location{{File: "file2.go", Line: 15}}},
			{ID: "id1", Str: "str1_modified", Locations: []Location{{File: "file2.go", Line: 25}}},
		},
	}

	file3 := &File{
		Name:     "file3",
		Header:   Header{},
		Nplurals: 1,
		Entries: Entries{
			{ID: "id4", Str: "str4", Locations: []Location{{File: "file3.go", Line: 30}}},
		},
	}

	mergedFile := MergeFiles(file1, file2, file3)

	expectedName := "file1_file2_file3"
	if mergedFile.Name != expectedName {
		t.Errorf("Expected merged file name to be %s, got %s", expectedName, mergedFile.Name)
	}

	if mergedFile.Nplurals != file1.Nplurals {
		t.Errorf("Expected Nplurals to be %d, got %d", file1.Nplurals, mergedFile.Nplurals)
	}

	expectedEntries := Entries{
		{
			ID:        "id1",
			Str:       "str1",
			Locations: []Location{{File: "file1.go", Line: 10}, {File: "file2.go", Line: 25}},
		},
		{ID: "id2", Str: "str2", Locations: []Location{{File: "file1.go", Line: 20}}},
		{ID: "id3", Str: "str3", Locations: []Location{{File: "file2.go", Line: 15}}},
		{ID: "id4", Str: "str4", Locations: []Location{{File: "file3.go", Line: 30}}},
	}

	if len(mergedFile.Entries) != len(expectedEntries) {
		t.Errorf("Expected %d entries, got %d", len(expectedEntries), len(mergedFile.Entries))
	}

	for i, entry := range mergedFile.Entries {
		if entry.ID != expectedEntries[i].ID || entry.Str != expectedEntries[i].Str {
			t.Errorf("Entry %d mismatch: expected %v, got %v", i, expectedEntries[i], entry)
		}
	}
}
