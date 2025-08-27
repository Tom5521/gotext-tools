package po_test

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/Tom5521/gotext-tools/v2/internal/util"
	"github.com/Tom5521/gotext-tools/v2/pkg/po"
	"github.com/Tom5521/gotext-tools/v2/pkg/po/compile"
	"github.com/Tom5521/gotext-tools/v2/pkg/po/parse"
)

func TestMsgcatMergeWithMsgcat(t *testing.T) {
	msgcat, err := exec.LookPath("msgcat")
	if err != nil {
		t.Skip(err)
		return
	}

	tmpdir := t.TempDir()

	file1Path := filepath.Join(tmpdir, "file1.po")
	file2Path := filepath.Join(tmpdir, "file2.po")
	file3Path := filepath.Join(tmpdir, "file3.po")
	outpath := filepath.Join(tmpdir, "out.po")

	file1 := &po.File{
		Name: "file1",
		Entries: po.Entries{
			{
				ID:  "id1",
				Str: "str1",
			},
			{
				ID:     "apple",
				Plural: "apples",
				Plurals: po.PluralEntries{
					{0, "manzana"},
					{1, "manzanas"},
				},
			},
			{
				ID:     "phone",
				Plural: "phones",
				Plurals: po.PluralEntries{
					{0, "telefono"},
					{1, "telefonos"},
				},
			},
		},
	}

	file2 := &po.File{
		Name: "file2",
		Entries: po.Entries{
			{
				ID: "id1",
			},
			{
				ID: "apple 2",
			},
			{
				ID:  "hello world",
				Str: "hola mundo",
			},
			{
				ID:  "meow",
				Str: "lol",
			},
		},
	}

	file3 := &po.File{
		Name: "file3",
		Entries: po.Entries{
			{
				ID:       "entry",
				Str:      "entrada",
				Comments: []string{"comment 1"},
			},
		},
	}

	// Write input
	if err = compile.PoToFile(file1, file1Path); err != nil {
		t.Error(err)
		return
	}
	if err = compile.PoToFile(file2, file2Path); err != nil {
		t.Error(err)
		return
	}
	if err = compile.PoToFile(file3, file3Path); err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name      string
		cmdArgs   []string
		mergeOpts []po.MsgcatMergeOption
	}{
		{
			name: "default",
		},
		{
			name:    "use first",
			cmdArgs: []string{"--use-first"},
			mergeOpts: []po.MsgcatMergeOption{
				po.MsgcatMergeWithUseFirst(true),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.cmdArgs = append(test.cmdArgs, file1Path, file2Path, file3Path)
			test.cmdArgs = append(test.cmdArgs, "-o", outpath)

			{
				var stderr bytes.Buffer
				cmd := exec.Command(msgcat, test.cmdArgs...)
				cmd.Stderr = &stderr
				if err = cmd.Run(); err != nil {
					t.Error(stderr.String())
					return
				}
			}

			expected, err := parse.Po(outpath)
			if err != nil {
				t.Error(err)
				return
			}

			obtained := po.MsgcatMergeFiles(
				[]*po.File{file1, file2, file3},
				test.mergeOpts...,
			)

			if !util.Equal(expected.Entries, obtained) {
				t.Error("obtained and expected differ!")
				fmt.Println(util.NamedDiff("obtained", "expected", expected.Entries, obtained))
				return
			}
		})
	}
}
