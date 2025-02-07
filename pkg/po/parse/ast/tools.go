package ast

import (
	"github.com/Tom5521/xgotext/internal/util"
)

func EqualNodeSlice(x, y []Node) bool {
	return util.Equal(x, y)
}

func EqualNodes(x, y Node) bool {
	return util.Equal(x, y)
}

func FormatNode(nodes ...Node) string {
	return util.Format(nodes...)
}
