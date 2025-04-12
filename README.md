# gotext-tools

[![Go Report Card](https://goreportcard.com/badge/github.com/Tom5521/xgotext)](https://goreportcard.com/report/github.com/Tom5521/xgotext)

A Go library and CLI toolkit for working with Gettext `.po` and `.mo` files. This project provides cross-platform tools for extracting, merging, and compiling translation files, with a focus on Go projects.

---

## Features

✅ **Cross-platform** – Written in Go, works on Windows, Linux, OpenBSD, NetBSD, and MacOS.  
✅ **Go-native support** – Designed with Go projects in mind.  
✅ **Easy to Install** – Pre-built binaries available, or build from source.

---

## CLI Tools

Located in `cli/`, these tools provide command-line utilities for managing `.po` files.

### `msgomerge`

A cross-platform alternative to `msgmerge`, used for updating `.po` files with new translations while preserving existing ones.

**Usage:**

```sh
msgomerge [old.po] [new.po] -o [output.po]
```

[More information here](/cli/msgomerge/README.md)

### `xgotext`

A Go-compatible version of `xgettext` for extracting translatable strings from Go source code.

**Usage:**

```sh
xgotext -o [output.po] [file1.go] [file2.go] ...
```

[More information here](/cli/xgotext/README.md)

---

📌 **Coming Soon:** More CLI tools for advanced Gettext operations.

---

## Library

The core library is located in `pkg/` and provides structured handling of `.po` and `.mo` files.

### `go/parse`

Extracts Gettext-compatible strings from Go source code. Useful for generating translation templates.

### `po`

The main package for working with `.po` files. Includes:

- **`Entry` & `Entries`** – Structured representation of translation entries.
- **`File`**
- **Sorting & Comparison** – Easily organize and compare translations.

### `po/compile`

Compiles parsed `.po` files into `.mo` (binary) or updated `.po` files.

### `po/parse`

Parsers for reading `.po` and `.mo` files into structured Go objects.

---

## Installation

### From Source

```sh
git clone https://github.com/Tom5521/gotext-tools
cd gotext-tools
go build ./cli/msgomerge
go build ./cli/xgotext
```

### Pre-built Binaries

Check the [Releases](https://github.com/Tom5521/gotext-tools/releases) page for pre-compiled executables.

---

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

---

## License

This project is open-source under the **MIT License**. See [LICENSE](https://github.com/Tom5521/gotext-tools/blob/main/LICENSE) for details.
