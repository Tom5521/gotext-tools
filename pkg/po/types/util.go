package types

import (
	"github.com/Tom5521/xgotext/internal/util"
)

func EqualLocations(l1, l2 []Location) bool {
	return util.Equal(l1, l2)
}

func EqualLocation(l1, l2 Location) bool {
	return util.Equal(l1, l2)
}

func EqualTranslation(t1, t2 Entry) bool {
	return util.Equal(t1, t2)
}

func EqualEntries(t1, t2 []Entry) bool {
	return util.Equal(t1, t2)
}
