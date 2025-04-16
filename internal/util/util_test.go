package util_test

import (
	"strings"
	"testing"

	"github.com/Tom5521/gotext-tools/v2/internal/util"
)

func TestFindLine(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error(r)
		}
	}()

	input := `a,b,c,d,e



f
g
h
i`

	expectedLine := 6
	index := strings.IndexRune(input, 'g')

	for i := 0; i < 2; i++ {
		var line int

		switch i {
		case 0:
			line = util.FindLine(input, index)
		case 1:
			line = util.FindLine([]byte(input), index)
		case 2:
			line = util.FindLine([]rune(input), index)
		}

		if line != expectedLine {
			t.Error("Unexpected line:")
			t.Error("Expected:", expectedLine)
			t.Error("Got:", line)
			break
		}
	}
}
