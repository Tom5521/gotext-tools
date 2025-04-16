//go:build ignore
// +build ignore

package util_test

import (
	"testing"

	"github.com/Tom5521/gotext-tools/v2/internal/util"
)

func TestPJWHash(t *testing.T) {
	if util.PJWHash("Hello!") != 0x04ec3311 {
		t.Fail()
	}
}
