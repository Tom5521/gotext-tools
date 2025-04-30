# msgofmt

A command-line tool for compiling Uniforum style `.po` files into binary `.mo` files. The tool converts human-readable PO files into machine-readable MO files, which are used by gettext for efficient runtime translation lookup.

## Features

- Compiles PO files into binary MO files
- Supports multiple input directories for file search
- Configurable output file name and location
- Option to force overwrite existing files
- Control over endianness of output file
- Option to exclude hash table from binary output
- Reads from standard input when input file is "-"

## Installation

```bash
curl -L -o $(go env GOPATH)/bin/msgofmt https://github.com/Tom5521/gotext-tools/releases/latest/download/msgofmt-$(go env GOOS)-$(go env GOARCH) && chmod +x $(go env GOPATH)/bin/msgofmt
```

## Usage

Basic usage:

```bash
msgofmt [flags] filename.po
```

### Command Line Options

- **Input/Output Options:**

  - `--output-file`, `-o`: Write output to specified file (default: "messages.mo").
  - `--directory`, `-D`: Add DIRECTORY to list for input files search.
  - `--force`, `-f`: Overwrite generated files if they already exist.

- **Binary Format Options:**

  - `--endianness`: Write out 32-bit numbers in the given byte order (options: "big", "little", or "native"; default: "native").
  - `--no-hash`: Binary file will not include the hash table.

- **Help:**
  - `--help`, `-h`: Display help information.

### Examples

Basic compilation with default output name:

```bash
msgofmt translations.po
```

Compile with custom output file name:

```bash
msgofmt -o my-translations.mo translations.po
```

Compile from standard input to standard output:

```bash
msgofmt -o - < translations.po
```

Compile with forced overwrite:

```bash
msgofmt -f -o output.mo translations.po
```

Compile with big-endian byte order:

```bash
msgofmt --endianness=big -o output.mo translations.po
```

Compile without hash table:

```bash
msgofmt --no-hash -o output.mo translations.po
```

Search for input file in additional directories:

```bash
msgofmt -D locale -D translations -o output.mo es.po
```

## How It Works

1. **Input File:**

   - Reads a PO file containing translations in human-readable format.
   - If input file is "-", reads from standard input.

2. **Compilation Process:**

   - Parses the PO file structure and validates its contents.
   - Converts text-based translations into binary format.
   - Optionally includes a hash table for faster lookups (unless disabled).
   - Handles byte order according to specified endianness.

3. **Output:**
   - Generates a binary MO file that can be used by gettext-enabled applications.
   - Outputs to specified file or standard output if "-" is specified.

## Acknowledgments

- [gettext](https://www.gnu.org/software/gettext/) - The GNU internationalization and localization system that defined the PO/MO file formats.
- [gotext](https://github.com/leonelquinteros/gotext) - The Go internationalization library this tool is designed to work with.
