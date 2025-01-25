package parse

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

func EqualNodeSlice(x, y []Node) bool {
	return slices.EqualFunc(x, y, EqualNodes)
}

func EqualNodes(x, y Node) bool {
	if x == y {
		return true
	}

	type1, type2 := reflect.TypeOf(x), reflect.TypeOf(y)
	value1, value2 := reflect.ValueOf(x), reflect.ValueOf(y)

	if type1.Kind() != type2.Kind() {
		return false
	}

	var equal bool
	for _, field := range reflect.VisibleFields(type1) {
		if !field.IsExported() {
			continue
		}
		equal = value1.FieldByIndex(field.Index).Interface() ==
			value2.FieldByIndex(field.Index).Interface()
		if !equal {
			break
		}
	}

	return equal
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
