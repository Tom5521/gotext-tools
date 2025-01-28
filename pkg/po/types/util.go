package types

import (
	"reflect"
	"slices"

	"github.com/Tom5521/xgotext/internal/util"
)

func CleanDuplicates(translations []Entry) (cleaned []Entry) {
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

func EqualLocations(l1, l2 []Location) bool {
	return slices.EqualFunc(l1, l2, EqualLocation)
}

func EqualLocation(l1, l2 Location) bool {
	return util.EqualFields(l1, l2)
}

func EqualTranslation(t1, t2 Entry) bool {
	return util.EqualFields(t1, t2)
}

func EqualEntries(t1, t2 []Entry) bool {
	return reflect.DeepEqual(t1, t2)
}
