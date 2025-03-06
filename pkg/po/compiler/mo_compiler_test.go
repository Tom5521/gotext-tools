package compiler_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/Tom5521/xgotext/pkg/po"
	"github.com/Tom5521/xgotext/pkg/po/compiler"
	"github.com/kr/pretty"
)

func TestMoCompiler(t *testing.T) {
	e := po.Entries{
		{ID: "id1", Str: "HELLO"},
		{ID: "id2", Str: "Hello2"},
		{ID: "id3", Str: "Hello3"},
	}
	c := compiler.NewMo(&po.File{Entries: e})

	var buf bytes.Buffer
	file, err := os.OpenFile("test.mo", os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		t.Error(err)
		return
	}

	err = c.ToWriter(io.MultiWriter(&buf, file))
	if err != nil {
		t.Error(err)
		return
	}

	cmd := exec.Command("msgunfmt", "test.mo")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Error(err)
		fmt.Println(string(out))
		pretty.Println("BYTES:\n", buf.Bytes())
		return
	}
}
