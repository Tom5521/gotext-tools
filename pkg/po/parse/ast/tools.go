package ast

import (
	"slices"

	"github.com/Tom5521/xgotext/internal/util"
)

func EqualNodeSlice(x, y []Node) bool {
	return slices.EqualFunc(x, y, EqualNodes)
}

func EqualNodes(x, y Node) bool {
	return util.EqualFields(x, y)
}

func FormatNode(nodes ...Node) string {
	return util.Format(nodes...)
}
