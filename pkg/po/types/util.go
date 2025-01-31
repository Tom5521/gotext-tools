package types

import (
	"reflect"
	"slices"

	"github.com/Tom5521/xgotext/internal/util"
)

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
