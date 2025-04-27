package compile_test

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/Tom5521/gotext-tools/v2/internal/util"
	"github.com/Tom5521/gotext-tools/v2/pkg/po"
	"github.com/Tom5521/gotext-tools/v2/pkg/po/compile"
	"github.com/Tom5521/gotext-tools/v2/pkg/po/parse"
)

func TestMoWithMsgfmt(t *testing.T) {
	msgunfmtPath, err := exec.LookPath("msgunfmt")
	if err != nil {
		t.Error(err)
		return
	}

	tmpDir := t.TempDir()
	outFile := filepath.Join(tmpDir, "out.mo")

	input := po.Entries{
		{Context: "My context :3", ID: "id1", Str: "HELLO"},
		{
			ID:     "id2",
			Plural: "helooows",
			Plurals: po.PluralEntries{
				po.PluralEntry{ID: 0, Str: "Holanda"},
				po.PluralEntry{ID: 1, Str: "Holandas"},
			},
		},
		{ID: "id3", Str: "Hello3"},
	}.SortFunc(po.CompareEntryByID)

	tests := []struct {
		name string
		opts []compile.MoOption
	}{
		{
			"Normal",
			nil,
		},
		{
			"With hash table",
			[]compile.MoOption{compile.MoWithHashTable(true)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.opts = append(test.opts, compile.MoWithForce(true))

			err = compile.MoToFile(input, outFile, test.opts...)
			if err != nil {
				t.Error(err)
				return
			}

			var stderr, stdout bytes.Buffer
			cmd := exec.Command(msgunfmtPath, outFile)
			cmd.Stderr = &stderr
			cmd.Stdout = &stdout

			err = cmd.Run()
			if err != nil {
				t.Error(stderr.String())
				return
			}

			parsed, err := parse.PoFromReader(&stdout, "test.po", parse.PoWithSkipHeader(true))
			if err != nil {
				t.Error(err)
				return
			}

			if !util.Equal(parsed.Entries, input) {
				t.Error(err)
			}
		})
	}
}
