package parse_test

import (
	"fmt"
	"testing"

	"github.com/Tom5521/gotext-tools/v2/internal/util"
	"github.com/Tom5521/gotext-tools/v2/pkg/go/parse"
	"github.com/Tom5521/gotext-tools/v2/pkg/po"
)

func TestParse(t *testing.T) {
	const input = `package main
import "github.com/leonelquinteros/gotext"

func main(){
	gotext.Get("Hello World!")
}`

	expected := po.Entries{
		{
			ID: "Hello World!",
			Locations: []po.Location{
				{
					Line: 5,
					File: "test.go",
				},
			},
		},
	}
	parser, err := parse.NewParserFromString(input, "test.go", parse.WithNoHeader(true))
	if err != nil {
		t.Error(err)
	}

	file := parser.Parse()
	if err = parser.Error(); err != nil {
		t.Error(err)
	}

	if !util.Equal(file.Entries, expected) {
		t.Error("expected and parsed differ!")
		fmt.Println(util.NamedDiff("expected", "parsed", expected, file.Entries))
		t.Fail()
		return
	}
}

func TestExtractAll(t *testing.T) {
	const input = `package main

import "github.com/leonelquinteros/gotext"

func main(){
	_ = "Hello World"
	a := "Hi world"
	b := "I love onions!"
	
	var eggs string = "sugar"
}`

	parser, err := parse.NewParserFromString(
		input,
		"test.go",
		parse.WithExtractAll(true),
		parse.WithNoHeader(true),
	)
	if err != nil {
		t.Error(err)
	}

	file := parser.Parse()
	if len(parser.Errors()) > 0 {
		t.Error(parser.Errors()[0])
	}

	expected := po.Entries{
		{
			ID: "Hello World",
			Locations: []po.Location{
				{
					File: "test.go",
					Line: 6,
				},
			},
		},
		{
			ID: "Hi world",
			Locations: []po.Location{
				{
					File: "test.go",
					Line: 7,
				},
			},
		},
		{
			ID: "I love onions!",
			Locations: []po.Location{
				{
					File: "test.go",
					Line: 8,
				},
			},
		},
		{
			ID: "sugar",
			Locations: []po.Location{
				{
					File: "test.go",
					Line: 10,
				},
			},
		},
	}

	if !util.Equal(file.Entries, expected) {
		t.Error("expected and parsed differ!")
		fmt.Println(util.NamedDiff("expected", "parsed", expected, file.Entries))
		return
	}
}

func TestExtractAll2(t *testing.T) {
	const input = `package main

import "github.com/leonelquinteros/gotext"

func main(){
	_ = "Hello World"
	a := "Hi world"
	b := "I love onions!"
	
	var eggs string = "sugar"
	const asadasd = "asad"

	for i := range "Hello World From a Loop!"{
		fmt.Println("Hello!")
	}
	for i := "hello";i != "from";i += "another loop"{

	}

	a := struct{x string}{"Hello from a struct!"}

	if "Hello From an if" != "Bye from an if"{
		print("no")
	}

	a := make(map[string]string)
	a["Hello from a key"] = "Hello from a value"
}`

	parser, _ := parse.NewParserFromString(
		input,
		"lol",
		parse.WithExtractAll(true),
		parse.WithNoHeader(true),
	)

	file := parser.Parse()
	if err := parser.Error(); err != nil {
		t.Error(err)
		return
	}

	expectedIDs := []string{
		"Hello World",
		"Hi world",
		"I love onions!",
		"sugar",
		"asad",
		"Hello World From a Loop!",
		"Hello!",
		"hello",
		"from",
		"another loop",
		"Hello from a struct!",
		"Hello From an if",
		"Bye from an if",
		"no",
		"Hello from a key",
		"Hello from a value",
	}

	ids := make([]string, len(file.Entries))
	for i, e := range file.Entries {
		ids[i] = e.ID
	}

	if !util.Equal(expectedIDs, ids) {
		t.Error("expected and parsed differ!")
		fmt.Println(util.NamedDiff("expected", "parsed", expectedIDs, ids))
		return
	}
}
