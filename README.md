# xgotext

A command-line tool for extracting translatable strings from Go source code that uses the [gotext](https://github.com/leonelquinteros/gotext) library for internationalization. This tool generates POT (Portable Object Template) files that can be used with standard translation tools.

## Features

- Extracts translatable strings from Go source files
- Supports all gotext translation functions (Get, GetD, GetN, GetC, GetND, GetNC, GetNDC)
- Handles multi-line strings
- Preserves context and plural forms
- Automatically removes duplicates while maintaining source references
- Supports excluding specific files or directories
- Provides verbose output option for debugging

## Installation

```bash
go install github.com/Tom5521/xgotext@latest
```

## Usage

Basic usage:

```bash
xgotext --input ./path/to/source --output messages.pot
```

### Command Line Options

- `--input`: Input file or directory path (default: ".")
- `--output`: Output POT file path (default: "default.pot")
- `--project-version`: Project version to include in the POT file
- `--lang`: Language code to include in the POT file
- `--nplurals`: Number of plurals for the target language (default: 2)
- `--exclude`: Files or directories to exclude from processing (can be specified multiple times)
- `--verbose`: Enable verbose output

### Examples

Extract strings from current directory:

```bash
xgotext --output messages.pot
```

Process a specific directory with version information:

```bash
xgotext --input ./cmd --output messages.pot --project-version "1.0.0"
```

Exclude certain directories:

```bash
xgotext --input . --output messages.pot --exclude vendor --exclude tests
```

Enable verbose output:

```bash
xgotext --input . --output messages.pot --verbose
```

## Supported Translation Functions

The tool extracts strings from the following gotext function calls:

- `gotext.Get(message)`
- `gotext.GetD(domain, message)`
- `gotext.GetN(message, plural, n)`
- `gotext.GetC(message, context)`
- `gotext.GetND(domain, message, plural, n)`
- `gotext.GetNC(message, plural, n, context)`
- `gotext.GetNDC(domain, message, plural, n, context)`

## Output Format

The generated POT file follows the standard gettext format, including:

- Source file references with line numbers
- Context information when available
- Plural forms
- Support for multi-line strings
- Empty msgstr fields ready for translation

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[MIT Licence](LICENCE)

## Acknowledgments

- [xgotext](https://github.com/leonelquinteros/gotext/tree/master/cli/xgotext) - The reason I did this, their tool does not work.
- [gotext](https://github.com/leonelquinteros/gotext) - The Go internationalization library this tool is designed to work with
- [spf13/pflag](https://github.com/spf13/pflag) - Command line flag parsing
