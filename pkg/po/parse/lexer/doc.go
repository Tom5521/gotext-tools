// Package lexer provides functionality for tokenizing PO (Portable Object) files.
// It parses the input content of a PO file into a sequence of tokens,
// which can then be used by a parser to construct an Abstract Syntax Tree (AST).
//
// Key Features:
// - Tokenizes PO file content into meaningful tokens such as comments,
// message IDs, message strings, contexts, and plural forms.
// - Handles escape characters within strings.
// - Supports keywords like `msgid`, `msgstr`, `msgctxt`, and `msgid_plural`.
//
// Example Usage:
//
//	input := `msgid "Hello"
//	msgstr "Hola"`
//	l := lexer.NewFromString(input)
//	for {
//	    tok := l.NextToken()
//	    if tok.Type == token.EOF {
//	        break
//	    }
//	    fmt.Println(tok)
//	}
//
// For more details, refer to the individual function and type documentation.
package lexer
