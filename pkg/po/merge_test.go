package po_test

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po"
	"github.com/Tom5521/xgotext/pkg/po/compiler"
	"github.com/Tom5521/xgotext/pkg/po/parse"
	"github.com/kr/pretty"
)

func TestMergeWithMsgmerge(t *testing.T) {
	msgmerge, err := exec.LookPath("msgmerge")
	if err != nil {
		t.Skip(err)
		return
	}

	tmpDir := t.TempDir()

	defPath := filepath.Join(tmpDir, "def.po")
	refPath := filepath.Join(tmpDir, "ref.po")
	outPath := filepath.Join(tmpDir, "out.po")

	defStruct := &po.File{Entries: po.Entries{
		{
			Context: "ctx1",
			ID:      "id1",
			Plural:  "plural id1",
			Plurals: []po.PluralEntry{
				{ID: 0, Str: "translated singular id1"},
				{ID: 1, Str: "translated plural id1"},
			},
		},
		{
			ID:     "id2",
			Plural: "plural id2",
			Plurals: []po.PluralEntry{
				{ID: 0, Str: "translated singular id2"},
				{ID: 1, Str: "translated plural id2"},
			},
		},
		{
			ID:  "id3",
			Str: "translated id3",
		},
		{
			ID:  "this must be removed.",
			Str: "old translation",
		},
	}}

	refStruct := &po.File{Entries: po.Entries{
		{
			Context: "ctx1",
			ID:      "id1",
			Plural:  "plural id1",
		},
		{
			ID:     "id2",
			Plural: "plural id2",
		},
		{
			ID: "id3",
		},
		{
			ID:     "id4",
			Plural: "plural id4",
		},
	}}

	// Write input.
	{
		comp := compiler.NewPo(defStruct)
		err = comp.ToFile(defPath)
		if err != nil {
			t.Error(err)
			return
		}
		comp.SetFile(refStruct)
		err = comp.ToFile(refPath)
		if err != nil {
			t.Error(err)
			return
		}
	}

	// Run command.
	{
		var stderr bytes.Buffer
		cmd := exec.Command(msgmerge, defPath, refPath, "-o", outPath)
		cmd.Stderr = &stderr
		if err = cmd.Run(); err != nil {
			t.Error(stderr.String())
			return
		}
	}

	parser, err := parse.NewPo(
		outPath,
		parse.PoWithSkipHeader(true),
		parse.PoWithIgnoreComments(true),
	)
	if err != nil {
		t.Error(err)
		return
	}

	expected := parser.Parse()
	if err = parser.Error(); err != nil {
		t.Error(err)
		return
	}

	getted := po.Merge(defStruct.Entries, refStruct.Entries).CleanObsoletes()

	if !util.Equal(expected.Entries, getted) {
		x, y := formatFileOrEntries(getted), formatFileOrEntries(expected)

		fmt.Println("DIFF:")
		for _, d := range pretty.Diff(expected.Entries, getted) {
			fmt.Println(d)
		}

		fmt.Println("Getted:\n", x)
		fmt.Println("Expected:\n", y)

		t.Fail()
		return
	}
}

func formatFileOrEntries[X *po.File | po.Entries](a X) string {
	var f *po.File
	switch v := any(a).(type) {
	case po.Entries:
		f = &po.File{Entries: v}
	case *po.File:
		f = v
	}

	return compiler.NewPo(f, compiler.PoWithOmitHeader(true)).ToString()
}
