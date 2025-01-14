package main

import (
	"strings"
	"unicode"
)

// cleanWhitespaces removes all whitespace characters from a string.
func cleanWhitespaces(s string) string {
	var builder strings.Builder
	builder.Grow(len(s)) // Pre-allocate capacity

	for _, char := range s {
		if !unicode.IsSpace(char) {
			builder.WriteRune(char)
		}
	}

	return builder.String()
}

// isEmpty returns true if the string contains only whitespace characters.
func isEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// findLine returns the line number for a given index in the string.
// Returns -1 if the index is out of bounds.
func findLine(str string, index int) int {
	if index < 0 || index >= len(str) {
		return -1
	}

	return 1 + strings.Count(str[:index], "\n")
}

// strContent extracts the content between the first pair of unescaped quotes.
// Returns an empty string if no valid quoted content is found.
func strContent(str string) string {
	var start, end int
	var foundStart, foundEnd bool
	for index, char := range str {
		if char == '"' {
			if !foundStart {
				start = index + 1
				foundStart = true
			} else if !foundEnd {
				if str[index-1] == '\\' {
					continue
				}
				end = index
				foundEnd = true
			} else {
				break
			}
		}
	}

	return str[start:end]
}
