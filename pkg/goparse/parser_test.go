package goparse_test

import (
	"testing"

	"github.com/Tom5521/xgotext/pkg/goparse"
	"github.com/Tom5521/xgotext/pkg/po/config"
	"github.com/Tom5521/xgotext/pkg/po/entry"
)

func TestParse(t *testing.T) {
	const input = `package main
import "github.com/leonelquinteros/gotext"

func main(){
	gotext.Get("Hello World!")
}`

	expected := []entry.Translation{
		{
			ID: "Hello World!",
			Locations: []entry.Location{
				{
					Line: 5,
					File: "test.go",
				},
			},
		},
	}
	cfg := config.Default()
	parser, err := goparse.NewParserFromBytes([]byte(input), "test.go", cfg, nil)
	if err != nil {
		t.Error(err)
	}

	translations, errs := parser.Parse()
	if len(errs) > 0 {
		t.Error(errs[0])
	}

	if !entry.CompareTranslations(translations, expected) {
		t.Log("Unexpected translations slice")
		t.Log("got:", translations)
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
	cfg := config.Default()
	cfg.ExtractAll = true
	parser, err := goparse.NewParserFromBytes([]byte(input), "test.go", cfg, nil)
	if err != nil {
		t.Error(err)
	}

	translations, errs := parser.Parse()
	if len(errs) > 0 {
		t.Error(errs[0])
	}

	expected := []entry.Translation{
		{
			ID: "Hello World",
			Locations: []entry.Location{
				{
					File: "test.go",
					Line: 6,
				},
			},
		},
		{
			ID: "Hi world",
			Locations: []entry.Location{
				{
					File: "test.go",
					Line: 7,
				},
			},
		},
		{
			ID: "I love onions!",
			Locations: []entry.Location{
				{
					File: "test.go",
					Line: 8,
				},
			},
		},
		{
			ID: "sugar",
			Locations: []entry.Location{
				{
					File: "test.go",
					Line: 10,
				},
			},
		},
	}

	if !entry.CompareTranslations(translations, expected) {
		t.Error("Unexpected translation")
		t.Log("got:", translations)
		t.Log("expected:", expected)
	}
}
