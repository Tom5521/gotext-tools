package entry

import (
	"slices"
)

func CompareLocations(l1, l2 []Location) bool {
	return slices.EqualFunc(l1, l2, CompareLocation)
}

func CompareLocation(l1, l2 Location) bool {
	return l1.File == l2.File && l1.Line == l2.Line
}

func CompareTranslation(t1, t2 Translation) bool {
	return (t1.ID == t2.ID && t1.Context == t2.Context && t1.Plural == t2.Plural) &&
		CompareLocations(t1.Locations, t2.Locations)
}

func CompareTranslations(t1, t2 []Translation) bool {
	return slices.EqualFunc(t1, t2, CompareTranslation)
}
