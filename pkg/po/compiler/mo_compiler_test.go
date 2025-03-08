package compiler_test

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po"
	"github.com/Tom5521/xgotext/pkg/po/compiler"
	"github.com/Tom5521/xgotext/pkg/po/parse"
	"github.com/kr/pretty"
)

func TestMoCompiler(t *testing.T) {
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
	}
	c := compiler.NewMo(&po.File{Entries: input})

	var buf, stderr, stdout bytes.Buffer
	cmd := exec.Command("msgunfmt", "-")
	cmd.Stdin = &buf
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

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

	parser, err := parse.NewPoFromReader(&stdout, "stdin")
	if err != nil {
		t.Error(err)
		return
	}

	parsed := parser.Parse().Entries
	if len(parser.Errors()) > 0 {
		t.Error(parser.Errors()[0])
		return
	}
	if !util.Equal(parsed, input) {
		t.Error("Sended and parsed differ!")
		t.Logf("SENDED:\n%v", input)
		t.Logf("PARSED:\n%v", parsed)
		t.Log("DIFF:")

		for _, d := range pretty.Diff(parsed, input) {
			fmt.Println(d)
		}
		return
	}
}
