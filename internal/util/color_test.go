package util_test

import (
	"testing"

	"github.com/Tom5521/gotext-tools/v2/internal/util"
)

func TestColor(t *testing.T) {
	t.Log(util.BgRed.Sprint(util.Black.Sprint("hi")))
}
