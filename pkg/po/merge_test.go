package po_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/Tom5521/gotext-tools/internal/util"
	"github.com/Tom5521/gotext-tools/pkg/po"
	"github.com/Tom5521/gotext-tools/pkg/po/compile"
	"github.com/Tom5521/gotext-tools/pkg/po/parse"
	fuzzy "github.com/paul-mannino/go-fuzzywuzzy"
	"github.com/sergi/go-diff/diffmatchpatch"
)

var dmp = diffmatchpatch.New()

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
		{
			ID: "id6",
		},
	}}

	// Write input.
	if err = compile.PoToFile(defStruct, defPath); err != nil {
		t.Error(err)
		return
	}
	if err = compile.PoToFile(refStruct, refPath); err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name      string
		cmdArgs   []string
		mergeOpts []po.MergeOption
	}{
		{
			"FuzzyMatch",
			[]string{},
			[]po.MergeOption{},
		},
		{
			"NoFuzzyMatch",
			[]string{"-N"},
			[]po.MergeOption{po.MergeWithFuzzyMatch(false)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mergeOpts = append(test.mergeOpts, po.MergeWithSort(true))
			test.cmdArgs = append(test.cmdArgs, defPath, refPath, "-o", outPath)
			// Run command.
			{
				var stderr bytes.Buffer
				cmd := exec.Command(msgmerge, test.cmdArgs...)
				cmd.Stderr = &stderr
				if err = cmd.Run(); err != nil {
					t.Error(stderr.String())
					return
				}
			}

			outBytes, err := os.ReadFile(outPath)
			if err != nil {
				t.Error(err)
				return
			}

			expected, err := parse.PoFromBytes(outBytes, outPath,
				parse.PoWithSkipHeader(true),
				parse.PoWithIgnoreComments(true),
			)
			if err != nil {
				t.Error(err)
				return
			}

			getted := po.Merge(defStruct.Entries, refStruct.Entries, test.mergeOpts...)

			if !util.Equal(expected.Entries, getted) {
				x, y := formatFileOrEntries(getted), formatFileOrEntries(expected)
				fmt.Println("--- STRUCT DIFF:")
				diff := dmp.DiffMain(util.Format(getted), util.Format(expected.Entries), false)
				fmt.Println(dmp.DiffPrettyText(diff))
				ratio := fuzzy.Ratio(x, y)
				fmt.Println("--- COMPILED MATCH RATIO:", ratio)
				diff = dmp.DiffMain(x, y, false)
				fmt.Println(dmp.DiffPrettyText(diff))

				t.Fail()
				return
			}
		})
	}
}

func formatFileOrEntries[X *po.File | po.Entries](a X) string {
	return compile.PoToString(a, compile.PoWithOmitHeader(true))
}
