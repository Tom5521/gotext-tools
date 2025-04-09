//go:build ignore
// +build ignore

package util

func PJWHash(s string) uint {
	var h uint
	var high uint

	for _, c := range s {
		h = (h << 4) + uint(c)
		high = h & 0xF0000000
		if high != 0 {
			h ^= high >> 24
		}
		h &= ^high
	}

	return h
}
