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
- Customizable headers and metadata for POT files
- Supports joining messages with existing POT files
- Configurable output directories and file naming

## Installation

```bash
go install github.com/Tom5521/gotext-tools/cli/xgotext@latest
```

## Usage

Basic usage:

```bash
xgotext [flags] file1 file2 file3...
```

### Command Line Options

- **Input/Output Options:**

  - `--files-from`, `-f`: Get list of input files from FILE.
  - `--directory`, `-D`: Add DIRECTORY to list for input files search. If input file is `-`, standard input is read.
  - `--output`, `-o`: Write output to specified file.
  - `--output-dir`, `-p`: Output files will be placed in directory DIR. If output file is `-`, output is written to standard output.
  - `--default-domain`, `-d`: Use NAME.pot for output (instead of messages.pot).

- **Parser Options:**

  - `--exclude`, `-X`: Specifies which files will be omitted.
  - `--extract-all`, `-a`: Extract all strings.
  - `--exclude-file`, `-x`: Entries from file are not extracted. File should be a PO or POT file.
  - `--join-existing`, `-j`: Join messages with existing file.

- **Header Options:**

  - `--lang`, `-l`: Language code to include in the POT file (default: "en").
  - `--nplurals`: Specify the number of plural forms of the language in question (default: 2).
  - `--msgid-bugs-address`: Set the reporting address for msgid bugs.
  - `--title`: Set the title of the POT file (default: "SOME DESCRIPTIVE TITLE").
  - `--copyright-holder`: Set the copyright holder in the output.
  - `--foreign-user`: Omit FSF copyright in output (equivalent to `--copyright-holder=''`).

- **Compiler Options:**
  - `--force-po`: Always write an output file even if no message is defined.
  - `--no-location`, `-n`: Do not write `#: filename:line` lines.
  - `--add-location`: Generate `#: filename:line` lines (default: "full"). Options: `full`, `file`, `never`.
  - `--omit-header`: Don’t write header with `msgid ""` entry.
  - `--package-name`: Set the package name in the header of the output.
  - `--package-version`: Set the package version in the header of the output.
  - `--msgstr-prefix`, `-m`: Use string as prefix for msgstr values.
  - `--msgstr-suffix`, `-M`: Use string as suffix for msgstr values.

### Examples

Extract strings from specific files:

```bash
xgotext -o messages.pot file1.go file2.go
```

Extract strings from all `.go` files in the current directory:

```bash
xgotext -o messages.pot *.go
```

Process a specific directory with version information:

```bash
xgotext -o messages.pot --package-version "1.0.0" ./cmd/*.go
```

Exclude certain files or directories:

```bash
xgotext -o messages.pot --exclude vendor/* --exclude tests/* *.go
```

Join messages with an existing POT file:

```bash
xgotext -o messages.pot --join-existing file1.go file2.go
```

Customize headers and metadata:

```bash
xgotext -o messages.pot --title "My Project" --copyright-holder "My Company" --msgid-bugs-address "bugs@example.com" file1.go file2.go
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
- Customizable headers and metadata

## Acknowledgments

- [xgotext](https://github.com/leonelquinteros/gotext/tree/master/cli/xgotext) - The reason I did this. ~,their tool does not work.~
- [gotext](https://github.com/leonelquinteros/gotext) - The Go internationalization library this tool is designed to work with
