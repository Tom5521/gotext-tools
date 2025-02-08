// Package ast provides the Abstract Syntax Tree (AST) for parsing and manipulating PO (Portable Object) files.
// It includes types and functions for parsing PO files into a structured AST,
// normalizing entries, merging multiple
// PO files, and handling various PO file constructs such as comments, message contexts, message IDs, and plural forms.
//
// The package is designed to be used in conjunction with the lexer and token packages to parse PO files into a
// structured representation that can be programmatically manipulated and queried.
//
// Key Features:
// - Parsing PO files into a structured AST.
// - Normalizing PO entries to ensure consistency.
// - Merging multiple PO files into a single file.
// - Handling various PO constructs like comments, message contexts, message IDs, and plural forms.
//
// Example Usage:
//
//	parser := ast.NewParser([]byte(input), "test.po")
//	normalizer, errs := parser.Normalizer()
//	if len(errs) > 0 {
//	    // Handle errors
//	}
//	normalizer.Normalize()
//	entries := normalizer.Entries()
//
//	// Merge multiple PO files
//	baseFile := &ast.File{Name: "base.po", Nodes: entries}
//	otherFile := &ast.File{Name: "other.po", Nodes: otherEntries}
//	ast.MergeFiles(baseFile, otherFile)
//
//	// Access the merged nodes
//	mergedNodes := baseFile.Nodes
//
// For more detailed usage, refer to the individual function and type documentation.
package ast
