// TODO: Remove repetitive code.
package parse_test

import (
	"reflect"
	"testing"

	. "github.com/Tom5521/xgotext/pkg/po/parse"
)

func TestGeneralComment(t *testing.T) {
	tok := Token{
		Type:    COMMENT,
		Literal: "# Hello World!",
	}

	expected := GeneralComment{
		Text: "Hello World!",
	}

	n, _ := (&Parser{}).Comment(tok)
	comment, ok := n.(GeneralComment)
	if !ok {
		t.Error("Unexpected node type:")
		t.Error("Expected type:", reflect.TypeOf(expected))
		t.Error("Got:", reflect.TypeOf(comment))
		return
	}

	if comment.Text != expected.Text {
		t.Error("Unexpected comment string:")
		t.Error("Expected string:", expected.Text)
		t.Error("Got:", comment.Text)
	}
}

func TestParseFlagComment(t *testing.T) {
	tok := Token{
		Type:    COMMENT,
		Literal: "#, my flag 123",
	}

	expected := FlagComment{
		Flag: "my flag 123",
	}

	n, _ := (&Parser{}).Comment(tok)
	flag, ok := n.(FlagComment)
	if !ok {
		t.Error("Unexpected node type:")
		t.Error("Expected type:", reflect.TypeOf(expected))
		t.Error("Got:", reflect.TypeOf(flag))
		return
	}

	if flag.Flag != expected.Flag {
		t.Error("Unexpected flag value:")
		t.Error("Expected:", expected.Flag)
		t.Error("Got:", flag.Flag)
	}
}

func TestParseLocationComment(t *testing.T) {
	tok := Token{
		Type:    COMMENT,
		Literal: "#: file:123",
	}

	expected := LocationComment{
		Line: 123,
		File: "file",
	}

	n, err := (&Parser{}).Comment(tok)
	if err != nil {
		t.Error(err)
		return
	}
	loc, ok := n.(LocationComment)
	if !ok {
		t.Error("Unexpected node type:")
		t.Error("Expected type:", reflect.TypeOf(expected))
		t.Error("Got:", reflect.TypeOf(loc))
		return
	}

	if loc.File != expected.File {
		t.Error("Unexpected file string:")
		t.Error("Expected string:", expected.File)
		t.Error("Got:", loc.File)
	}
	if loc.Line != expected.Line {
		t.Error("Unexpected line integer:")
		t.Error("Expected integer:", expected.Line)
		t.Error("Got:", loc.Line)
	}
}

func TestParseMsgid(t *testing.T) {
	input := `msgid ""
"Hello World"`

	expected := Msgid{
		ID: "\nHello World",
	}

	p := NewParserFromString(input, "test.po")
	errs := p.Parse()

	if len(errs) > 0 {
		t.Error("Unexpected error:", errs[0])
		return
	}

	if len(p.Nodes()) == 0 {
		t.Error("Unexpected amount of nodes (0)")
		t.Error("Want: 1")
		return
	}

	msgid := p.Nodes()[0].(Msgid)
	if msgid.ID != expected.ID {
		t.Error("Unexpected msgid id:")
		t.Error("Expected:", expected.ID)
		t.Error("Got:", msgid.ID)
	}
}

func TestParseMsgstr(t *testing.T) {
	input := `msgstr ""
"Hello World"`

	expected := Msgstr{Str: "\nHello World"}

	p := NewParserFromString(input, "test.po")
	errs := p.Parse()

	if len(errs) > 0 {
		t.Error("Unexpected error:", errs[0])
		return
	}

	if len(p.Nodes()) == 0 {
		t.Error("Unexpected amount of nodes (0)")
		t.Error("Want: 1")
		return
	}

	msgstr := p.Nodes()[0].(Msgstr)
	if msgstr.Str != expected.Str {
		t.Error("Unexpected msgstr text:")
		t.Error("Expected:", expected.Str)
		t.Error("Got:", msgstr.Str)
	}
}
