package compiler_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/Tom5521/xgotext/pkg/po"
	"github.com/Tom5521/xgotext/pkg/po/compiler"
	"github.com/kr/pretty"
)

func TestMoCompiler(t *testing.T) {
	e := po.Entries{
		{Context: "My context :3", ID: "id1", Str: "HELLO"},
		{ID: "id2", Str: "Hello2", Plural: "helooows", Plurals: po.PluralEntries{po.PluralEntry{ID: 0, Str: "Holanda"}, po.PluralEntry{ID: 1, Str: "Holandas"}}},
		{ID: "id3", Str: "Hello3"},
	}
	c := compiler.NewMo(&po.File{Entries: e})

	var buf, stderr bytes.Buffer
	cmd := exec.Command("msgunfmt", "-")
	cmd.Stdin = &buf
	cmd.Stderr = &stderr
	cmd.Stdout = os.Stdout

	err := c.ToWriter(&buf)
	if err != nil {
		t.Error(err)
		return
	}

	err = cmd.Run()
	if err != nil {
		t.Error(err)
		fmt.Println(stderr.String())
		pretty.Println("BYTES:\n", buf.Bytes())
		return
	}
}
