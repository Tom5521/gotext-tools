package color_test

import (
	"testing"

	"github.com/Tom5521/gotext-tools/v2/internal/color"
)

func TestColor(t *testing.T) {
	t.Log(color.BgRed.Sprint(color.Black.Sprint("hi")))
}
