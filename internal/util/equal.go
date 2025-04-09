package util

import (
	"math"
	"reflect"

	fuzzy "github.com/paul-mannino/go-fuzzywuzzy"
)

func FuzzyEqual(x, y string) bool {
	return fuzzy.Ratio(x, y) >= 80
}

type visitedPairs map[[2]uintptr]struct{}

func Equal[X, Y any](x X, y Y) bool {
	return equal(reflect.ValueOf(x), reflect.ValueOf(y), make(visitedPairs))
}

func equal(v1, v2 reflect.Value, visited visitedPairs) bool {
	if v1.Kind() == reflect.Pointer || v1.Kind() == reflect.Interface {
		if v1.IsNil() || v2.IsNil() {
			return v1.IsNil() == v2.IsNil()
		}
		return equal(v1.Elem(), v2.Elem(), visited)
	}

	if v1.Type() != v2.Type() {
		return false
	}

	if v1.CanAddr() && v2.CanAddr() {
		addr1, addr2 := v1.UnsafeAddr(), v2.UnsafeAddr()
		pair := [2]uintptr{addr1, addr2}
		if _, found := visited[pair]; found {
			return true
		}
		visited[pair] = struct{}{}
	}

	switch v1.Kind() {
	case reflect.Bool:
		return v1.Bool() == v2.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v1.Int() == v2.Int()
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr:
		return v1.Uint() == v2.Uint()
	case reflect.Float32, reflect.Float64:
		return floatEqual(v1.Float(), v2.Float())
	case reflect.Complex64, reflect.Complex128:
		return complexEqual(v1.Complex(), v2.Complex())
	case reflect.String:
		return v1.String() == v2.String()
	case reflect.Array, reflect.Slice:
		return sliceEqual(v1, v2, visited)
	case reflect.Map:
		return mapEqual(v1, v2, visited)
	case reflect.Struct:
		return structEqual(v1, v2, visited)
	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		return false
	default:
		return false
	}
}

func floatEqual(a, b float64) bool {
	const epsilon = 1e-9
	return math.Abs(a-b) < epsilon
}

func complexEqual(a, b complex128) bool {
	return floatEqual(real(a), real(b)) && floatEqual(imag(a), imag(b))
}

func sliceEqual(v1, v2 reflect.Value, visited visitedPairs) bool {
	if v1.Len() != v2.Len() {
		return false
	}
	for i := 0; i < v1.Len(); i++ {
		if !equal(v1.Index(i), v2.Index(i), visited) {
			return false
		}
	}
	return true
}

func mapEqual(v1, v2 reflect.Value, visited visitedPairs) bool {
	if v1.Len() != v2.Len() {
		return false
	}
	for _, key := range v1.MapKeys() {
		val1 := v1.MapIndex(key)
		val2 := v2.MapIndex(key)
		if !val2.IsValid() || !equal(val1, val2, visited) {
			return false
		}
	}
	return true
}

func structEqual(v1, v2 reflect.Value, visited visitedPairs) bool {
	for i := 0; i < v1.NumField(); i++ {
		field := v1.Type().Field(i)
		if field.IsExported() {
			if !equal(v1.Field(i), v2.Field(i), visited) {
				return false
			}
		}
	}
	return true
}
