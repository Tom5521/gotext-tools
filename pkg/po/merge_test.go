package po_test

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po"
	"github.com/Tom5521/xgotext/pkg/po/compiler"
	"github.com/Tom5521/xgotext/pkg/po/parse"
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

	defStruct := &po.File{Entries: po.Entries{}}
	refStruct := &po.File{Entries: po.Entries{}}

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

	parser, err := parse.NewPo(outPath, parse.PoWithSkipHeader(true))
	if err != nil {
		t.Error(err)
		return
	}

	expected := parser.Parse()
	if err = parser.Error(); err != nil {
		t.Error(err)
		return
	}

	getted := refStruct.Merge(defStruct)

	if !util.Equal(expected.Entries, getted.Entries) {
		t.Fail()
		return
	}
}
