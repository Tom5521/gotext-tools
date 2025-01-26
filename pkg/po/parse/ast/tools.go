package ast

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/Tom5521/xgotext/internal/util"
)

func EqualNodeSlice(x, y []Node) bool {
	return slices.EqualFunc(x, y, EqualNodes)
}

func EqualNodes(x, y Node) bool {
	return util.EqualFields(x, y)
}

func FormatNode(nodes ...Node) string {
	var b strings.Builder

	for i, n := range nodes {
		value := reflect.ValueOf(n)
		typeOf := reflect.TypeOf(n)

		isPointer := typeOf.Kind() == reflect.Pointer

		// Get values
		if isPointer {
			value = value.Elem()
			typeOf = typeOf.Elem()

			fmt.Fprint(&b, "&")
		}

		fmt.Fprint(&b, typeOf.String())
		fmt.Fprintln(&b, "{")
		for _, field := range reflect.VisibleFields(typeOf) {
			if !field.IsExported() {
				continue
			}
			fmt.Fprintf(&b, "  %s: %#v,\n", field.Name, value.FieldByIndex(field.Index).Interface())
		}
		fmt.Fprint(&b, "}")
		if i != len(nodes)-1 {
			b.WriteByte(',')
		}
		b.WriteByte('\n')
	}

	return b.String()
}
