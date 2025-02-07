package ast_test

import (
	"testing"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po/parse/ast"
	"github.com/kr/pretty"
)

func diff(a, b any) string {
	var str string

	for _, d := range pretty.Diff(a, b) {
		str += d + "\n"
	}

	return str
}

func TestEntries_SortByLine(t *testing.T) {
	tests := []struct {
		name     string
		entries  ast.Entries
		expected ast.Entries
	}{
		{
			name: "Basic sorting",
			entries: ast.Entries{
				{LocationComments: []*ast.LocationComment{{Line: 3}}},
				{LocationComments: []*ast.LocationComment{{Line: 1}}},
				{LocationComments: []*ast.LocationComment{{Line: 2}}},
			},
			expected: ast.Entries{
				{LocationComments: []*ast.LocationComment{{Line: 1}}},
				{LocationComments: []*ast.LocationComment{{Line: 2}}},
				{LocationComments: []*ast.LocationComment{{Line: 3}}},
			},
		},
		{
			name: "Entries without LocationComments",
			entries: ast.Entries{
				{},
				{LocationComments: []*ast.LocationComment{{Line: 1}}},
			},
			expected: ast.Entries{
				{LocationComments: []*ast.LocationComment{{Line: 1}}},
				{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sorted := tt.entries.SortByLine()
			if !util.Equal(sorted, tt.expected) {
				t.Errorf("expected %v, got %v", pretty.Sprint(tt.expected), pretty.Sprint(sorted))
			}
		})
	}
}

func TestEntries_CleanDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		entries  ast.Entries
		expected ast.Entries
	}{
		{
			name: "Basic duplicates",
			entries: ast.Entries{
				{Msgid: &ast.Msgid{ID: "A"}},
				{Msgid: &ast.Msgid{ID: "B"}},
				{Msgid: &ast.Msgid{ID: "A"}},
			},
			expected: ast.Entries{
				{Msgid: &ast.Msgid{ID: "A"}},
				{Msgid: &ast.Msgid{ID: "B"}},
			},
		},
		{
			name: "No duplicates",
			entries: ast.Entries{
				{Msgid: &ast.Msgid{ID: "A"}},
				{Msgid: &ast.Msgid{ID: "B"}},
			},
			expected: ast.Entries{
				{Msgid: &ast.Msgid{ID: "A"}},
				{Msgid: &ast.Msgid{ID: "B"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleaned := tt.entries.CleanDuplicates()
			if !util.Equal(cleaned, tt.expected) {
				t.Errorf("expected %v, got %v", pretty.Sprint(tt.expected), pretty.Sprint(cleaned))
			}
		})
	}
}

func TestEntries_Sort(t *testing.T) {
	tests := []struct {
		name     string
		entries  ast.Entries
		expected ast.Entries
	}{
		{
			name: "Sort by file and line",
			entries: ast.Entries{
				{LocationComments: []*ast.LocationComment{{File: "b/file.go", Line: 3}}},
				{LocationComments: []*ast.LocationComment{{File: "a/file.go", Line: 1}}},
				{LocationComments: []*ast.LocationComment{{File: "a/file.go", Line: 2}}},
			},
			expected: ast.Entries{
				{LocationComments: []*ast.LocationComment{{File: "a/file.go", Line: 1}}},
				{LocationComments: []*ast.LocationComment{{File: "a/file.go", Line: 2}}},
				{LocationComments: []*ast.LocationComment{{File: "b/file.go", Line: 3}}},
			},
		},
		{
			name: "Files without LocationComments",
			entries: ast.Entries{
				{},
				{LocationComments: []*ast.LocationComment{{File: "a/file.go", Line: 1}}},
				{},
			},
			expected: ast.Entries{
				{},
				{},
				{LocationComments: []*ast.LocationComment{{File: "a/file.go", Line: 1}}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sorted := tt.entries.Sort()
			if !util.Equal(sorted, tt.expected) {
				t.Errorf("expected %v, got %v", pretty.Sprint(tt.expected), pretty.Sprint(sorted))
				t.Error(diff(tt.expected, sorted))
			}
		})
	}
}

func TestMergeFiles(t *testing.T) {
	tests := []struct {
		name     string
		base     *ast.File
		files    []*ast.File
		expected string
	}{
		{
			name: "Basic merge",
			base: &ast.File{Name: "base", Nodes: []ast.Node{ast.Entry{Msgid: &ast.Msgid{ID: "A"}}}},
			files: []*ast.File{
				{Name: "file1", Nodes: []ast.Node{ast.Entry{Msgid: &ast.Msgid{ID: "B"}}}},
				{Name: "file2", Nodes: []ast.Node{ast.Entry{Msgid: &ast.Msgid{ID: "A"}}}},
			},
			expected: "base_file1_file2",
		},
		{
			name: "No additional files",
			base: &ast.File{
				Name:  "base",
				Nodes: []ast.Node{ast.Entry{Msgid: &ast.Msgid{ID: "A"}}},
			},
			files:    []*ast.File{},
			expected: "base",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast.MergeFiles(tt.base, tt.files...)
			if tt.base.Name != tt.expected {
				t.Errorf("expected file name %s, got %s", tt.expected, tt.base.Name)
			}
		})
	}
}

func TestEntries_Solve(t *testing.T) {
	tests := []struct {
		name     string
		entries  ast.Entries
		expected ast.Entries
	}{
		{
			name: "Resolve duplicates with Msgstr",
			entries: ast.Entries{
				{Msgid: &ast.Msgid{ID: "A"}, Msgstr: &ast.Msgstr{Str: "Hello"}},
				{Msgid: &ast.Msgid{ID: "A"}, Msgstr: nil},
				{Msgid: &ast.Msgid{ID: "B"}, Msgstr: &ast.Msgstr{Str: "World"}},
			},
			expected: ast.Entries{
				{Msgid: &ast.Msgid{ID: "A"}, Msgstr: &ast.Msgstr{Str: "Hello"}},
				{Msgid: &ast.Msgid{ID: "B"}, Msgstr: &ast.Msgstr{Str: "World"}},
			},
		},
		{
			name: "Resolve duplicates with Plurals",
			entries: ast.Entries{
				{Msgid: &ast.Msgid{ID: "A"}, Plurals: []*ast.MsgstrPlural{{Str: "One"}}},
				{Msgid: &ast.Msgid{ID: "A"}, Plurals: []*ast.MsgstrPlural{}},
				{Msgid: &ast.Msgid{ID: "B"}, Plurals: []*ast.MsgstrPlural{{Str: "Two"}}},
			},
			expected: ast.Entries{
				{Msgid: &ast.Msgid{ID: "A"}, Plurals: []*ast.MsgstrPlural{{Str: "One"}}},
				{Msgid: &ast.Msgid{ID: "B"}, Plurals: []*ast.MsgstrPlural{{Str: "Two"}}},
			},
		},
		{
			name: "Resolve duplicates without Msgstr or Plurals",
			entries: ast.Entries{
				{Msgid: &ast.Msgid{ID: "A"}},
				{Msgid: &ast.Msgid{ID: "A"}},
				{Msgid: &ast.Msgid{ID: "B"}},
			},
			expected: ast.Entries{
				{Msgid: &ast.Msgid{ID: "A"}},
				{Msgid: &ast.Msgid{ID: "B"}},
			},
		},
		{
			name: "Resolve duplicates with LocationComments",
			entries: ast.Entries{
				{
					Msgid:            &ast.Msgid{ID: "A"},
					LocationComments: []*ast.LocationComment{{File: "file1.go", Line: 1}},
				},
				{
					Msgid:            &ast.Msgid{ID: "A"},
					LocationComments: []*ast.LocationComment{{File: "file2.go", Line: 2}},
				},
			},
			expected: ast.Entries{
				{
					Msgid: &ast.Msgid{ID: "A"},
					LocationComments: []*ast.LocationComment{
						{File: "file1.go", Line: 1},
						{File: "file2.go", Line: 2},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			solved := tt.entries.Solve()
			if !util.Equal(solved, tt.expected) {
				t.Errorf("expected %v, got %v", pretty.Sprint(tt.expected), pretty.Sprint(solved))
			}
		})
	}
}
