package goparse_test

import (
	"slices"
	"testing"

	"github.com/Tom5521/xgotext/pkg/goparse"
	"github.com/Tom5521/xgotext/pkg/poconfig"
	"github.com/Tom5521/xgotext/pkg/poentry"
)

func TestParse(t *testing.T) {
	const input = `package main
import "github.com/leonelquinteros/gotext"

func main(){
	gotext.Get("Hello World!")
}`

	expected := []poentry.Translation{
		{
			ID: "Hello World!",
			Locations: []poentry.Location{
				{
					Line: 5,
					File: "test.go",
				},
			},
		},
	}
	cfg := poconfig.DefaultConfig()
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
		func(e1 poentry.Translation, e2 poentry.Translation) bool {
			return (e1.ID == e2.ID && e1.Context == e2.Context && e1.Plural == e2.Plural) &&
				slices.EqualFunc(
					e1.Locations,
					e2.Locations,
					func(e1 poentry.Location, e2 poentry.Location) bool {
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
