package po_test

import (
	"testing"

	"github.com/Tom5521/xgotext/pkg/po"
)

func TestMergeFiles(t *testing.T) {
	file1 := &po.File{
		Name: "file1",
		Entries: po.Entries{
			{ID: "id1", Str: "str1", Locations: []po.Location{{File: "file1.go", Line: 10}}},
			{ID: "id2", Str: "str2", Locations: []po.Location{{File: "file1.go", Line: 20}}},
		},
	}

	file2 := &po.File{
		Name: "file2",
		Entries: po.Entries{
			{ID: "id3", Str: "str3", Locations: []po.Location{{File: "file2.go", Line: 15}}},
			{
				ID:        "id1",
				Str:       "str1_modified",
				Locations: []po.Location{{File: "file2.go", Line: 25}},
			},
		},
	}

	file3 := &po.File{
		Name: "file3",
		Entries: po.Entries{
			{ID: "id4", Str: "str4", Locations: []po.Location{{File: "file3.go", Line: 30}}},
		},
	}

	mergedFile := file1.MergeWithOptions([]*po.File{file2, file3}, po.MergeWithFuzzyMatch(false))

	expectedName := "file1_file2_file3"
	if mergedFile.Name != expectedName {
		t.Errorf("Expected merged file name to be %s, got %s", expectedName, mergedFile.Name)
	}

	expectedEntries := po.Entries{
		{
			ID:        "id1",
			Str:       "str1",
			Locations: []po.Location{{File: "file1.go", Line: 10}, {File: "file2.go", Line: 25}},
		},
		{ID: "id2", Str: "str2", Locations: []po.Location{{File: "file1.go", Line: 20}}},
		{ID: "id3", Str: "str3", Locations: []po.Location{{File: "file2.go", Line: 15}}},
		{ID: "id4", Str: "str4", Locations: []po.Location{{File: "file3.go", Line: 30}}},
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
