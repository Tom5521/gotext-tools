# xgotext

A command-line tool that extracts message IDs from Go source files for internationalization purposes. The tool scans Go files for calls to the `gotext.Get()` function and generates a POT (Portable Object Template) file that can be used with translation tools.

## Features

- Extracts message IDs from single files or entire directories
- Generates standard POT files with file references and line numbers
- Supports file exclusion patterns
- Handles custom import aliases for the gotext package
- Preserves original message formatting and context

## Installation

```bash
go install github.com/Tom5521/xgotext@latest
```

## Usage

Basic usage:

```bash
xgotext --input ./path/to/source --output messages.pot
```

### Available Flags

- `--input`: Input file or directory path (default: ".")
- `--output`: Output POT file path (default: "default.pot")
- `--proyect-version`: Project version to include in the POT file
- `--lang`: Language code to include in the POT file
- `--exclude`: Files or directories to exclude from processing (can be specified multiple times)
- `--verbose`: Enable verbose output

### Examples

Extract messages from a specific directory:

```bash
xgotext --input ./src --output messages.pot
```

Extract messages with project metadata:

```bash
xgotext --input . --output messages.pot --proyect-version "1.0.0" --lang "en-US"
```

Exclude specific files or directories:

```bash
xgotext --input . --output messages.pot --exclude vendor --exclude "**/*_test.go"
```

## How It Works

1. The tool recursively scans Go source files in the specified directory
2. It looks for imports of the `github.com/leonelquinteros/gotext` package
3. For files that import gotext, it extracts all strings passed to the `Get()` function
4. Generates a POT file containing all found message IDs with their source locations

## Generated POT File Format

The tool generates a standard POT file with entries in the following format:

```
#: path/to/file.go:123
msgid "Original message"
msgstr ""
```

## Requirements

- Go 1.21 or higher
- github.com/leonelquinteros/gotext
- github.com/spf13/pflag
- github.com/gookit/color

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request to [github.com/Tom5521/xgotext](https://github.com/Tom5521/xgotext).

## License

This is licensed under the MIT license.

[LICENCE]
