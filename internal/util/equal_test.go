package util_test

import (
	"testing"

	"github.com/Tom5521/gotext-tools/internal/util"
)

func TestEqual(t *testing.T) {
	type A struct {
		X *A
	}

	var a, b A
	a.X = &a

	b.X = &b

	// Stack overflow.
	util.Equal(a, b)
}
