package compile_test

import (
	"strings"
	"testing"

	"github.com/Tom5521/gotext-tools/v2/pkg/po"
	"github.com/Tom5521/gotext-tools/v2/pkg/po/compile"
)

func TestCompile(t *testing.T) {
	var builder strings.Builder

	eb := compile.EntryBuilder{
		Entry: po.Entry{
			Flags:             []string{"plural"},
			Comments:          []string{"Plural forms for items"},
			ExtractedComments: []string{"Shopping cart module"},
			Previous:          []string{},
			Obsolete:          false,
			ID:                "%d item\nlol\nline2",
			Context:           "shopping_cart",
			Plural:            "%d items",
			Plurals: po.PluralEntries{
				{ID: 0, Str: "%d artículo"},
				{ID: 1, Str: "%d artículos"},
			},
			Locations: po.Locations{
				{Line: 88, File: "cart.go"},
				{Line: 92, File: "cart.go"},
			},
		},
		Builder: &builder,
		Config: compile.DefaultPoConfig(
			compile.PoWithWordWrap(true),
		),
	}

	eb.Build()
	t.Log(builder.String())
}
