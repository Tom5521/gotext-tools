//go:build !go1.21
// +build !go1.21

package slices

import "sort"

type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

func SortFunc[S ~[]E, E any](x S, cmp func(a, b E) int) {
	sort.Slice(x, func(i, j int) bool {
		return cmp(x[i], x[j]) < 0
	})
}
