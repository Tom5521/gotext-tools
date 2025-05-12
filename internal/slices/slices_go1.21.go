//go:build go1.21
// +build go1.21

package slices

import (
	"cmp"
	"slices"
)

type Ordered = cmp.Ordered

func SortFunc[S ~[]E, E any](x S, cmp func(a, b E) int) {
	slices.SortFunc(x, cmp)
}
