# msgocat

A command-line tool for concatenating and merging Uniforum style `.po` files. The tool combines multiple PO files, finds common messages across them, and generates a unified PO file with merged translations, comments, and file positions.

## Features

- Concatenates and merges multiple PO files
- Finds messages common to two or more PO files
- Filters messages based on occurrence frequency
- Cumulates translations, comments, and file positions
- Supports sorting output by file location or message content
- Configurable location comment generation
- Color output support
- Reads from standard input when input file is "-"

## Installation

```bash
curl -L -o $(go env GOPATH)/bin/msgocat https://github.com/Tom5521/gotext-tools/releases/latest/download/msgocat-$(go env GOOS)-$(go env GOARCH) && chmod +x $(go env GOPATH)/bin/msgocat
```

## Usage

Basic usage:

```bash
msgocat [flags] [file1.po file2.po ...]
```

### Command Line Options

- **Input/Output Options:**
  - `--output-file`, `-o`: Write output to specified file (default: "-" for standard output).
  - `--directory`, `-D`: Add DIRECTORY to list for input files search.
  - `--files-from`, `-f`: Get list of input files from specified file.

- **Message Filtering Options:**
  - `--more-than`, `->`: Print messages with more than this many definitions (default: 0).
  - `--less-than`, `-<`: Print messages with less than this many definitions (default: max uint64).
  - `--unique`, `-u`: Shorthand for `--less-than=2`, requests that only unique messages be printed.
  - `--use-first`: Use first available translation for each message, don't merge several translations.

- **Output Formatting Options:**
  - `--add-location`: Generate '#: filename:line' lines (default: "full").
  - `--no-location`: Do not write '#: filename:line' lines.
  - `--no-wrap`: Do not break long message lines.
  - `--sort-output`, `-s`: Generate sorted output.
  - `--sort-by-file`, `-F`: Sort output by file location.
  - `--color`: Use colors and other text attributes (options: "always", "never", "auto"; default: "auto").
  - `--lang`: Set 'Language' field in the header entry.

- **Help:**
  - `--help`, `-h`: Display help information.

### Aliases

The tool can also be invoked as:

- `msgocat`
- `msgcat`
- `cat`

### Examples

Basic concatenation of multiple PO files:

```bash
msgocat file1.po file2.po file3.po
```

Concatenate and output to file:

```bash
msgocat -o merged.po file1.po file2.po
```

Find messages common to at least 3 files:

```bash
msgocat --more-than=3 *.po
```

Find unique messages (appearing in only one file):

```bash
msgocat --unique *.po
```

Use first available translation and suppress location comments:

```bash
msgocat --use-first --no-location *.po
```

Sort output by file location:

```bash
msgocat --sort-by-file *.po
```

Read input files from a list:

```bash
msgocat --files-from=filelist.txt
```

Search for input files in additional directories:

```bash
msgocat -D locale -D translations *.po
```

## How It Works

1. **Input Processing:**
   - Reads multiple PO files containing translations
   - If input file is "-", reads from standard input
   - Can read file lists from specified files using `--files-from`

2. **Message Analysis:**
   - Identifies messages common across multiple files
   - Applies frequency filters (`--more-than`, `--less-than`) to select messages
   - Cumulates translations, comments, and file positions from all sources

3. **Translation Selection:**
   - By default, merges translations from all files
   - With `--use-first`, uses the translation from the first file that defines each message

4. **Output Generation:**
   - Generates a unified PO file with selected messages
   - Formats location comments according to specified options
   - Can sort output by message content or file location
   - Outputs to specified file or standard output

## Acknowledgments

- [gettext](https://www.gnu.org/software/gettext/) - The GNU internationalization and localization system that defined the PO file format.
- [gotext](https://github.com/leonelquinteros/gotext) - The Go internationalization library this tool is designed to work with.
