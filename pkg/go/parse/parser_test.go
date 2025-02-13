package parse_test

import (
	"testing"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/go/parse"
	goparse "github.com/Tom5521/xgotext/pkg/go/parse"
	"github.com/Tom5521/xgotext/pkg/po/types"
	"github.com/kr/pretty"
)

func TestParse(t *testing.T) {
	const input = `package main
import "github.com/leonelquinteros/gotext"

func main(){
	gotext.Get("Hello World!")
}`

	expected := types.Entries{
		{
			ID: "Hello World!",
			Locations: []types.Location{
				{
					Line: 5,
					File: "test.go",
				},
			},
		},
	}
	parser, err := goparse.NewParserFromString(input, "test.go")
	if err != nil {
		t.Error(err)
	}

	file := parser.Parse()
	if len(parser.Errors()) > 0 {
		t.Error(parser.Errors()[0])
	}

	if !types.EqualEntries(file.Entries[1:], expected) {
		t.Log("Unexpected entries slice")
		t.Log("got:", file.Entries)
		t.Log("expected:", expected)
		t.FailNow()
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

	parser, err := goparse.NewParserFromString(
		input,
		"test.go",
		goparse.WithExtractAll(true),
	)
	if err != nil {
		t.Error(err)
	}

	file := parser.Parse()
	if len(parser.Errors()) > 0 {
		t.Error(parser.Errors()[0])
	}

	expected := types.Entries{
		{
			ID: "Hello World",
			Locations: []types.Location{
				{
					File: "test.go",
					Line: 6,
				},
			},
		},
		{
			ID: "Hi world",
			Locations: []types.Location{
				{
					File: "test.go",
					Line: 7,
				},
			},
		},
		{
			ID: "I love onions!",
			Locations: []types.Location{
				{
					File: "test.go",
					Line: 8,
				},
			},
		},
		{
			ID: "sugar",
			Locations: []types.Location{
				{
					File: "test.go",
					Line: 10,
				},
			},
		},
	}

	if !util.Equal(file.Entries[1:], expected) {
		t.Error("Unexpected translation")
		t.Log("got:", file.Entries)
		t.Log("expected:", expected)
		t.Log("DIFF:")
		for _, d := range pretty.Diff(file.Entries, expected) {
			t.Log(d)
		}
	}
}

func BenchmarkParse(b *testing.B) {
	const input = `package main
import "github.com/leonelquinteros/gotext"

func main(){
	gotext.Get("Hello World!")
}`

	for i := 0; i < b.N; i++ {
		parser, err := parse.NewParserFromString(input, "test.go")
		if err != nil {
			b.Error(err)
		}

		parser.Parse()
		if len(parser.Errors()) > 0 {
			b.Error(parser.Errors()[0])
		}
	}
}

func BenchmarkExtractAll(b *testing.B) {
	const input = `package main

import "github.com/leonelquinteros/gotext"

func main(){
	_ = "Hello World"
	a := "Hi world"
	b := "I love onions!"
	
	var eggs string = "sugar"
}`

	for i := 0; i < b.N; i++ {
		parser, err := parse.NewParserFromString(
			input,
			"test.go",
			parse.WithExtractAll(true),
		)
		if err != nil {
			b.Error(err)
		}

		parser.Parse()
		if len(parser.Errors()) > 0 {
			b.Error(parser.Errors()[0])
		}
	}
}
