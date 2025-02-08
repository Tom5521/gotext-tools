// Package parse provides tools for parsing Go source files to extract translations and handle various configurations.
// It is designed to work with the "gotext" library,
// allowing users to extract strings that need translation from Go code.
//
// The package includes functionality to:
// - Parse Go files and extract translation entries.
// - Handle different configurations for parsing, such as excluding specific paths or extracting all strings.
// - Process abstract syntax trees (AST) to identify translation function calls and string literals.
// - Generate translation entries with location information for use in PO files.
//
// The main components of the package are:
// - Parser: The core struct that manages the parsing process, including file handling and configuration.
// - File: Represents a Go source file being processed, containing methods to parse and extract translations.
// - Util: Provides utility functions for AST traversal and configuration validation.
// - Process: Contains logic for processing AST nodes and
// extracting translation entries from function calls and string literals.
//
// Example usage:
//
//	cfg := parsers.Config{
//	    ExtractAll: true, // Extract all strings, not just those marked for translation.
//	}
//	parser, err := parse.NewParserFromFiles([]string{"example.go"}, cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	file := parser.Parse()
//	for _, entry := range file.Entries {
//	    fmt.Println(entry.ID, entry.Locations)
//	}
//
// This package is part of the xgotext project,
// which aims to provide comprehensive tools for internationalization in Go.
package parse
