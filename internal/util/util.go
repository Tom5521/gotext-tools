package util

import (
	"bytes"
	"strings"

	"github.com/kr/pretty"
)

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

func ROverNumber[T Number](x T) chan T {
	ch := make(chan T)

	go func() {
		defer close(ch)
		for i := T(0); i < x; i++ {
			ch <- i
		}
	}()

	return ch
}

func SafeSliceAccess[T any](slice []T, index int) T {
	var value T

	if index >= 0 && index < len(slice) {
		value = slice[index]
	}

	return value
}

func Format[T any](args ...T) string {
	var a []any
	for _, arg := range args {
		a = append(a, arg)
	}

	return pretty.Sprint(a...)
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
