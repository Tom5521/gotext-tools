package util

import (
	"bytes"

	"github.com/Tom5521/xgotext/pkg/poentry"
)

// FindLine returns the line number for a given index in the slice of bytes.
// Returns -1 if the index is out of bounds.
func FindLine[T ~int](content []byte, index T) int {
	if index < 0 || int(index) >= len(content) {
		return -1
	}

	return 1 + bytes.Count(content[:index], []byte{'\n'})
}

func CleanDuplicates(translations []poentry.Translation) (cleaned []poentry.Translation) {
	seenID := make(map[string]int)

	for _, translation := range translations {
		idIndex, ok := seenID[translation.ID]
		if ok {
			if translation.Context == cleaned[idIndex].Context {
				cleaned[idIndex].Locations = append(
					cleaned[idIndex].Locations,
					translation.Locations...)
				continue
			}
		}
		seenID[translation.ID] = len(cleaned)
		cleaned = append(cleaned, translation)
	}

	return
}
