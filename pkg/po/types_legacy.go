//go:build !go1.23
// +build !go1.23

package po

type (
	EntriesOrFile interface{ Entries | *File | File }
	Cmp[X any]    func(a, b X) int
)
