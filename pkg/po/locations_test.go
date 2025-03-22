package po_test

import (
	"math/rand"
	"slices"
	"testing"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po"
)

func TestSortLocations(t *testing.T) {
	expected := po.Locations{
		{File: "A", Line: 1},
		{File: "A", Line: 2},
		{File: "B", Line: 1},
		{File: "B", Line: 10},
		{File: "C", Line: 1},
		{File: "C", Line: 2},
	}

	sorted := slices.Clone(expected)
	rand.Shuffle(len(expected), func(i, j int) {
		sorted[i], sorted[j] = sorted[j], sorted[i]
	})

	sorted = sorted.Sort()

	if !util.Equal(expected, sorted) {
		t.Fail()
	}
}
