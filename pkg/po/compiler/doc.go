// Package compiler provides functionality to compile translation entries into PO (Portable Object) files.
// It handles formatting, configuration, and output generation for PO files, making it easy to integrate with translation workflows.
//
// The package is divided into three main components:
//   - `compiler.go`: Core functionality for compiling translations.
//   - `config.go`: Configuration options for the compilation process.
//   - `format.go`: Formatting utilities for PO file headers and entries.
//
// # Overview
//
// The compiler takes a `types.File` containing translation entries and a `Config` object to customize the compilation process.
// It supports outputting the compiled translations to a file, string, or byte slice.
//
// # Example Usage
//
// The following example demonstrates how to use the compiler to generate a PO file:
//
//	package main
//
//	import (
//		"fmt"
//		"github.com/Tom5521/xgotext/pkg/po/types"
//		"github.com/Tom5521/xgotext/pkg/compiler"
//	)
//
//	func main() {
//		// Create a sample translation file.
//		file := &types.File{
//			Header: types.Header{
//				Fields: []types.HeaderField{
//					{Key: "Project-Id-Version", Value: "1.0"},
//				},
//			},
//			Entries: []types.Entry{
//				{
//					ID: "Hello, World!",
//					Locations: []types.Location{
//						{File: "main.go", Line: 10},
//					},
//				},
//			},
//		}
//
//		// Configure the compiler.
//		config := compiler.Config{
//			PackageName:     "myapp",
//			CopyrightHolder: "My Company",
//			Title:           "My App Translations",
//		}
//
//		// Compile the translations to a string.
//		compiler := compiler.Compiler{File: file, Config: config}
//		output := compiler.CompileToString()
//		fmt.Println(output)
//	}
//
// # Components
//
// ## Compiler
//
// The `Compiler` struct is the main entry point for the compilation process. It provides the following methods:
//   - `CompileToWriter`: Writes the compiled translations to an `io.Writer`.
//   - `CompileToFile`: Writes the compiled translations to a file.
//   - `CompileToString`: Returns the compiled translations as a string.
//   - `CompileToBytes`: Returns the compiled translations as a byte slice.
//
// ## Configuration
//
// The `Config` struct defines various options for the compilation process, such as:
//   - `ForcePo`: Force creation of the PO file, even if it exists.
//   - `OmitHeader`: Omit the header from the output.
//   - `PackageName`: Name of the package being translated.
//   - `CopyrightHolder`: Copyright holder for the translations.
//   - `ForeignUser`: Indicates if the file is for a foreign user (public domain).
//   - `Title`: Title of the PO file.
//   - `NoLocation`: Suppress location comments in the output.
//   - `AddLocation`: Control how location comments are added ("full", "file", or "never").
//   - `MsgstrPrefix`: Prefix to add to `msgstr` values.
//   - `MsgstrSuffix`: Suffix to add to `msgstr` values.
//
// The `Validate` method ensures that the configuration is valid and free of conflicts.
//
// ## Formatting
//
// The `format.go` file contains utilities for formatting the PO file header and individual translation entries. Key functions include:
//   - `formatHeader`: Formats the header of the PO file.
//   - `formatEntry`: Formats a single translation entry.
//   - `formatString`: Formats a string to be compatible with PO file syntax.
//   - `formatMultiline`: Formats a string as a PO-compatible multiline string.
//   - `fixSpecialChars`: Escapes special characters in a string.
//
// # Notes
//
// - The compiler ensures that the output conforms to the PO file format, making it compatible with tools like `gettext`.
// - The `CleanDuplicates` method is used to remove duplicate translation entries before compilation.
//
// For more details, refer to the individual source files.
package compiler
