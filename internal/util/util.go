package util

import (
	"bytes"
	"reflect"
	"strings"

	"github.com/Tom5521/xgotext/pkg/po/types"
)

func EqualFields(x, y any) bool {
	if x == y {
		return true
	}

	typeX, typeY := reflect.TypeOf(x), reflect.TypeOf(y)
	valueX, valueY := reflect.ValueOf(x), reflect.ValueOf(y)

	if typeX.Kind() == reflect.Pointer {
		typeX = typeX.Elem()
		valueX = valueX.Elem()
	}
	if typeY.Kind() == reflect.Pointer {
		typeY = typeY.Elem()
		valueY = valueY.Elem()
	}

	if typeX.Kind() != typeY.Kind() {
		return false
	}

	var equal bool
	for _, field := range reflect.VisibleFields(typeX) {
		if !field.IsExported() {
			continue
		}
		v1, v2 := valueX.FieldByIndex(field.Index), valueY.FieldByIndex(field.Index)
		if field.Type.Kind() == reflect.Struct {
			equal = EqualFields(v1.Interface(), v2.Interface())
		} else {
			equal = v1.Interface() == v2.Interface()
		}

		if !equal {
			break
		}
	}

	return equal
}

func CountRunes(slice []rune, target rune) int {
	count := 0
	for _, r := range slice {
		if r == target {
			count++
		}
	}
	return count
}

// FindLine returns the line number (1-based) for a given index within the content.
// The function accepts content as string, []rune, or []byte and determines the line number
// by counting newline characters ('\n') that appear before the specified index.
//
// Generic Parameters:
//   - T: integer numeric type for the index (constrained to integer types)
//   - B: content type, constrained to []rune, []byte, or string
//
// Parameters:
//   - content: the content to analyze (string, []rune, or []byte)
//   - index: position in the content for which to determine the line number
//
// Returns:
//   - int: line number (starting from 1) corresponding to the index
//   - returns -1 if the index is out of range (negative or greater than content length)
//
// Example usage:
//
//	text := "Line 1\nLine 2\nLine 3"
//	lineNum := FindLine(text, 10) // Returns 2, as index 10 is on the second line
//
// The function internally handles different content types:
//   - For strings: uses strings.Count to count '\n' occurrences
//   - For []rune: uses an auxiliary CountRunes function
//   - For []byte: uses bytes.Count to count '\n' occurrences
//
// Implementation notes:
//   - The function uses a type assertion to determine the content type at runtime
//   - Counts the number of newline characters ('\n') before the given index
//   - Returns line number starting from 1
//   - Returns -1 if the index is out of bounds (negative or exceeding content length)
//
// Edge cases:
//   - If the content is empty or the index is out of range, the function returns -1.
func FindLine[T ~int, B []rune | []byte | string](content B, index T) int {
	if index < 0 || int(index) >= len(content) {
		return -1
	}

	switch c := any(content).(type) {
	case string:
		return strings.Count(c[:index], "\n") + 1
	case []rune:
		return CountRunes(c[:index], '\n') + 1
	default:
		return bytes.Count(c.([]byte)[:index], []byte{'\n'}) + 1
	}
}

func CleanDuplicates(translations []types.Entry) (cleaned []types.Entry) {
	seenID := make(map[string]int)

	for _, translation := range translations {
		idIndex, ok := seenID[translation.ID]
		if ok {
			if translation.Context == cleaned[idIndex].Context {
				cleaned[idIndex].Locations = append(
					cleaned[idIndex].Locations,
					translation.Locations...)
				continue
			}
		}
		seenID[translation.ID] = len(cleaned)
		cleaned = append(cleaned, translation)
	}

	return
}
