// Package types defines the core data structures and utilities for handling PO (Portable Object) files.
// It includes types for managing headers, entries (translatable strings), locations in source code, and plural forms.
// The package also provides functions for sorting,
// merging, and normalizing entries, as well as generating and manipulating headers.
//
// Key Features:
// - **Header Management**: Supports creating, updating, and retrieving header fields for PO files.
// - **Entry Management**: Handles translatable strings,
// including their context, plural forms, and source code locations.
// - **Sorting and Merging**: Provides utilities for sorting entries by ID, file, line, and merging multiple PO files.
// - **Normalization**: Ensures consistency by removing duplicates and resolving conflicts in entries.
//
// Example Usage:
//
//	// Create a default header
//	header := types.DefaultHeader()
//	header.Set("Project-Id-Version", "MyProject 1.0")
//
//	// Create a new PO file
//	file := types.NewFile("example.po", types.Entry{ID: "greeting", Str: "Hello"})
//
//	// Merge two PO files
//	mergedFile := types.MergeFiles(file1, file2)
//
// For more details, refer to the individual types and functions.
package po
