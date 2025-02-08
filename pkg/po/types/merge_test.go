package types_test

import (
	"testing"

	"github.com/Tom5521/xgotext/pkg/po/types"
)

func TestMergeFiles(t *testing.T) {
	file1 := &types.File{
		Name: "file1",
		Entries: types.Entries{
			{ID: "id1", Str: "str1", Locations: []types.Location{{File: "file1.go", Line: 10}}},
			{ID: "id2", Str: "str2", Locations: []types.Location{{File: "file1.go", Line: 20}}},
		},
	}

	file2 := &types.File{
		Name: "file2",
		Entries: types.Entries{
			{ID: "id3", Str: "str3", Locations: []types.Location{{File: "file2.go", Line: 15}}},
			{
				ID:        "id1",
				Str:       "str1_modified",
				Locations: []types.Location{{File: "file2.go", Line: 25}},
			},
		},
	}

	file3 := &types.File{
		Name: "file3",
		Entries: types.Entries{
			{ID: "id4", Str: "str4", Locations: []types.Location{{File: "file3.go", Line: 30}}},
		},
	}

	mergedFile := *file1
	types.MergeFiles(&mergedFile, file2, file3)

	expectedName := "file1_file2_file3"
	if mergedFile.Name != expectedName {
		t.Errorf("Expected merged file name to be %s, got %s", expectedName, mergedFile.Name)
	}

	expectedEntries := types.Entries{
		{
			ID:        "id1",
			Str:       "str1",
			Locations: []types.Location{{File: "file1.go", Line: 10}, {File: "file2.go", Line: 25}},
		},
		{ID: "id2", Str: "str2", Locations: []types.Location{{File: "file1.go", Line: 20}}},
		{ID: "id3", Str: "str3", Locations: []types.Location{{File: "file2.go", Line: 15}}},
		{ID: "id4", Str: "str4", Locations: []types.Location{{File: "file3.go", Line: 30}}},
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
