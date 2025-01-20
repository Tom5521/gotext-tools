package goparse_test

import (
	"slices"
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
	cfg := config.DefaultConfig()
	parser, err := goparse.NewParserFromBytes([]byte(input), "test.go", cfg)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	translations, errs := parser.Parse()
	if len(errs) > 0 {
		t.Log(errs[0])
		t.FailNow()
	}

	if !slices.EqualFunc(
		translations,
		expected,
		func(e1 entry.Translation, e2 entry.Translation) bool {
			return (e1.ID == e2.ID && e1.Context == e2.Context && e1.Plural == e2.Plural) &&
				slices.EqualFunc(
					e1.Locations,
					e2.Locations,
					func(e1 entry.Location, e2 entry.Location) bool {
						return e1.File == e2.File && e1.Line == e2.Line
					},
				)
		},
	) {
		t.Log("Unexpected translations slice")
		t.Log(translations)
		t.Log("expected: ", expected)
		t.FailNow()
	}
}
