// Package parse provides tools for parsing Go source files to extract translations and handle various configurations.
// It is designed to work with the "gotext" library, allowing users to extract translatable strings from Go code.
//
// ### Key Features
//
// 1. **Translation Extraction**:
//   - Extracts strings meant for translation, including those from specific function calls (e.g., `gotext.Get`).
//   - Optionally extracts all string literals in the code (useful for broader analysis).
//
// 2. **Configurable Parsing**:
//   - Exclude specific paths or files from processing.
//   - Handle custom configurations such as headers, fuzzy matching, and verbose logging.
//
// 3. **AST-Based Processing**:
//   - Processes Go abstract syntax trees (AST) to identify translation function calls and string literals.
//   - Generates translation entries with metadata, including file location and line number.
//
// 4. **Integration with PO Files**:
//   - Produces translation entries compatible with PO (Portable Object) files, facilitating localization workflows.
//
// ### Main Components
//
// - **Parser**:
//   - Manages the parsing process, including file handling, configuration, and error reporting.
//   - Processes files and extracts translation entries.
//
// - **File**:
//   - Represents a single Go source file.
//   - Handles parsing, package detection, and translation extraction.
//
// - **Process**:
//   - Contains logic to traverse AST nodes and extract translation entries from function calls and string literals.
//
// - **Config**:
//   - Defines configurations for customizing the parsing process
//
// (e.g., verbosity, exclusion rules, header customization).
//
// ### Example Usage
//
// ```go
//
//	cfg := parse.Config{
//	    ExtractAll: true, // Extract all string literals, not just those marked for translation.
//	    Verbose:    true, // Enable verbose logging.
//	}
//
// parser, err := parse.NewParserFromFiles([]string{"example.go"}, parse.WithConfig(cfg))
//
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// file := parser.Parse()
//
//	for _, entry := range file.Entries {
//	    fmt.Printf("String: %s, Locations: %v\n", entry.ID, entry.Locations)
//	}
//
// ```
//
// ### About
// This package is part of the xgotext project, which provides comprehensive tools for internationalization in Go.
package parse
