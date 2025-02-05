package ast_test

import (
	"testing"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po/parse/ast"
	"github.com/kr/pretty"
)

func TestEntries_SortByLine(t *testing.T) {
	entries := ast.Entries{
		ast.Entry{LocationComments: []*ast.LocationComment{{Line: 3}}},
		ast.Entry{LocationComments: []*ast.LocationComment{{Line: 1}}},
		ast.Entry{LocationComments: []*ast.LocationComment{{Line: 2}}},
	}

	sorted := entries.SortByLine()
	expected := ast.Entries{
		ast.Entry{LocationComments: []*ast.LocationComment{{Line: 1}}},
		ast.Entry{LocationComments: []*ast.LocationComment{{Line: 2}}},
		ast.Entry{LocationComments: []*ast.LocationComment{{Line: 3}}},
	}

	if !util.Equal(sorted, expected) {
		t.Errorf("expected %v, got %v", pretty.Sprint(expected), pretty.Sprint(sorted))
	}
}

func TestEntries_CleanDuplicates(t *testing.T) {
	entries := ast.Entries{
		ast.Entry{Msgid: &ast.Msgid{ID: "A"}},
		ast.Entry{Msgid: &ast.Msgid{ID: "B"}},
		ast.Entry{Msgid: &ast.Msgid{ID: "A"}},
	}

	cleaned := entries.CleanDuplicates()

	expected := ast.Entries{
		ast.Entry{Msgid: &ast.Msgid{ID: "A"}},
		ast.Entry{Msgid: &ast.Msgid{ID: "B"}},
	}

	if !util.Equal(cleaned, expected) {
		t.Errorf("expected %v, got %v", pretty.Sprint(expected), pretty.Sprint(cleaned))
	}
}

func TestEntries_Sort(t *testing.T) {
	entries := ast.Entries{
		ast.Entry{LocationComments: []*ast.LocationComment{{File: "b/file.go", Line: 3}}},
		ast.Entry{LocationComments: []*ast.LocationComment{{File: "a/file.go", Line: 1}}},
		ast.Entry{LocationComments: []*ast.LocationComment{{File: "a/file.go", Line: 2}}},
	}

	sorted := entries.Sort()
	expected := ast.Entries{
		ast.Entry{LocationComments: []*ast.LocationComment{{File: "a/file.go", Line: 1}}},
		ast.Entry{LocationComments: []*ast.LocationComment{{File: "a/file.go", Line: 2}}},
		ast.Entry{LocationComments: []*ast.LocationComment{{File: "b/file.go", Line: 3}}},
	}

	if !util.Equal(sorted, expected) {
		t.Errorf("expected %v, got %v", expected, sorted)
	}
}

func TestMergeFiles(t *testing.T) {
	base := &ast.File{Name: "base", Nodes: []ast.Node{ast.Entry{Msgid: &ast.Msgid{ID: "A"}}}}
	file1 := &ast.File{Name: "file1", Nodes: []ast.Node{ast.Entry{Msgid: &ast.Msgid{ID: "B"}}}}
	file2 := &ast.File{Name: "file2", Nodes: []ast.Node{ast.Entry{Msgid: &ast.Msgid{ID: "A"}}}}

	ast.MergeFiles(base, file1, file2)

	expectedName := "base_file1_file2"
	if base.Name != expectedName {
		t.Errorf("expected file name %s, got %s", expectedName, base.Name)
	}
}
