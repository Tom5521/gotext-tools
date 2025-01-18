package parser

import (
	"go/token"
	"strings"
)

// findLine returns the line number for a given index in the string.
// Returns -1 if the index is out of bounds.
func findLine(str string, index token.Pos) int {
	if index < 0 || int(index) >= len(str) {
		return -1
	}

	return 1 + strings.Count(str[:index], "\n")
}

func cleanDuplicates(translations []Translation) (cleaned []Translation) {
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
