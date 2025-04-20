package util

// PJWHash computes a hash value for a string using the PJW (Elf) hash algorithm.
func PJWHash(str string) uint32 {
	var h, g uint32

	for _, c := range str {
		h = (h << 4) + uint32(c)
		g = h & 0xF0000000 // Check the top 4 bits
		if g != 0 {
			h ^= g >> 24 // XOR with the high bits
			h &= ^g      // Clear the top 4 bits
		}
	}
	return h
}
