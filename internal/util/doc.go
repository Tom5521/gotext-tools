// Package util provides a collection of utility functions designed to facilitate
// common operations such as deep equality comparison, string formatting, and
// content analysis. These functions are particularly useful for scenarios requiring
// dynamic type handling, reflection-based comparison, and text processing.
//
// Functions:
//
//   - Equal: Performs a deep comparison between two values of any type. It accounts for
//     different types of data structures, including pointers, arrays, slices, maps,
//     structs, and interfaces. The function uses reflection to handle complex cases,
//     ensuring accurate comparisons, even for floating-point and complex numbers.
//
//   - Format: Converts a set of input values into a formatted string representation.
//     This function supports generic types and leverages the pretty-printing library
//     to improve readability of structured data, making it useful for debugging
//     and logging purposes.
//
//   - FindLine: Determines the line number corresponding to a given index within a text
//     content. It accepts different types of content, including strings, rune slices,
//     and byte slices. The function efficiently counts newline characters ('\n')
//     preceding the given index to determine the correct line number. If the index
//     is out of range, it returns -1.
//
//   - CountRunes: Counts the occurrences of a specified rune within a slice of runes.
//     This function is particularly useful when analyzing text content where character
//     frequency is relevant, such as counting specific symbols or delimiters.
//
// Implementation Details:
//
// The package makes extensive use of reflection to perform type-agnostic operations,
// particularly in the `Equal` function. It ensures robust handling of nested data
// structures, taking into account exported struct fields, nil pointers, and different
// numeric types. Additionally, special care is taken when comparing floating-point
// numbers, using a small tolerance (epsilon) to mitigate precision errors.
//
// The `FindLine` function is optimized to efficiently scan through different types of
// textual content, leveraging specialized functions like `strings.Count` and `bytes.Count`
// for their respective input types. This ensures both correctness and performance,
// especially when processing large files or extensive textual data.
package util
