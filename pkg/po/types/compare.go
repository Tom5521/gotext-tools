package types

import (
	"slices"
)

func EqualLocations(l1, l2 []Location) bool {
	return slices.EqualFunc(l1, l2, EqualLocation)
}

func EqualLocation(l1, l2 Location) bool {
	return l1.File == l2.File && l1.Line == l2.Line
}

func EqualTranslation(t1, t2 Entry) bool {
	return (t1.ID == t2.ID && t1.Context == t2.Context && t1.Plural == t2.Plural) &&
		EqualLocations(t1.Locations, t2.Locations)
}

func EqualTranslations(t1, t2 []Entry) bool {
	return slices.EqualFunc(t1, t2, EqualTranslation)
}
