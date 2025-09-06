# msgounfmt

A command-line tool for converting binary `.mo` message catalogs back into human-readable Uniforum style `.po` files. The tool decompiles binary MO files into text-based PO files, reversing the compilation process performed by tools like msgofmt.

## Features

- Converts binary MO files to text-based PO files
- Supports multiple input files for batch processing
- Configurable output file name and location
- Option to force PO file generation even for empty catalogs
- Control over output formatting (line wrapping)
- Ability to generate sorted output
- Error handling options for non-critical data issues
- Reads from standard input when input file is "-"

## Installation

```bash
curl -L -o $(go env GOPATH)/bin/msgounfmt https://github.com/Tom5521/gotext-tools/releases/latest/download/msgounfmt-$(go env GOOS)-$(go env GOARCH) && chmod +x $(go env GOPATH)/bin/msgounfmt
```

## Usage

Basic usage:

```bash
msgounfmt [flags] [file1.mo file2.mo ...]
```

### Command Line Options

- **Input/Output Options:**
  - `--output`, `-o`: Write output to specified file (default: "-" for standard output).
  - If no input file is given or if it is "-", standard input is read.

- **Output Formatting Options:**
  - `--no-wrap`: Do not break long message lines into several lines.
  - `--sort-output`, `-s`: Generate sorted output.
  - `--force-po`, `-f`: Write PO file even if empty.
  - `--color`: Use colors and other text attributes (options: "always", "never", "auto"; default: "auto").

- **Error Handling Options:**
  - `--ignore-errors`: Skip non-critical errors in the data such as duplicate and/or unsorted entries.

- **Help:**
  - `--help`, `-h`: Display help information.

### Aliases

The tool can also be invoked as:

- `msgounfmt`
- `msgunfmt`
- `unfmt`

### Examples

Basic conversion of a single MO file:

```bash
msgounfmt messages.mo
```

Convert multiple MO files:

```bash
msgounfmt en.mo es.mo fr.mo
```

Convert with custom output file:

```bash
msgounfmt -o output.po messages.mo
```

Convert without line wrapping:

```bash
msgounfmt --no-wrap messages.mo
```

Generate sorted output:

```bash
msgounfmt --sort-output messages.mo
```

Force PO file creation even for empty catalogs:

```bash
msgounfmt --force-po empty.mo
```

Read from standard input and output to file:

```bash
cat messages.mo | msgounfmt -o output.po
```

Ignore non-critical errors in the data:

```bash
msgounfmt --ignore-errors corrupted.mo
```

## How It Works

1. **Input Processing:**
   - Reads binary MO files containing compiled translations
   - If no input file is specified or if "-" is given, reads from standard input
   - Supports processing multiple input files in a single operation

2. **Decompilation Process:**
   - Parses the binary MO file structure
   - Extracts message strings, translations, and metadata
   - Handles different byte orders and file formats
   - Validates file integrity while providing options to ignore non-critical errors

3. **Output Generation:**
   - Converts binary data back into text-based PO format
   - Formats messages with proper PO file structure
   - Optionally sorts messages for consistent output
   - Controls line wrapping based on user preference
   - Outputs to specified file or standard output

4. **Error Handling:**
   - Detects and reports critical errors in binary data
   - Optionally ignores non-critical issues like duplicate entries
   - Ensures valid PO file generation even with imperfect input

## Use Cases

- **Debugging Translations:** Convert compiled MO files back to PO format to inspect and debug translation issues
- **Recovery:** Recover lost PO files from deployed MO files
- **Migration:** Convert MO files from other systems for use with gotext-based applications
- **Analysis:** Inspect the contents of binary message catalogs

## Acknowledgments

- [gettext](https://www.gnu.org/software/gettext/) - The GNU internationalization and localization system that defined the MO/PO file formats.
- [gotext](https://github.com/leonelquinteros/gotext) - The Go internationalization library this tool is designed to work with.
