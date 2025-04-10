# msgomerge

A command-line tool for merging two Uniforum style `.po` files together. The tool combines translations from an existing PO file (`def.po`) with up-to-date source references from a reference POT file (`ref.pot`), preserving comments and applying fuzzy matching where exact matches are not found.

## Features

- Merges translations from an existing PO file with source references from a POT file
- Preserves comments from the existing PO file
- Supports fuzzy matching for improved translation reuse
- Customizable location tags for source references
- Option to disable fuzzy matching for strict merging
- Supports additional translation libraries (compendium files)
- Configurable output file and directory
- Language field customization in the header
- Option to update the existing PO file in-place

## Installation

```bash
go install github.com/Tom5521/gotext-tools/cli/msgomerge@latest
```

## Usage

Basic usage:

```bash
msgomerge [flags] def.po ref.pot
```

### Command Line Options

- **Input/Output Options:**

  - `--output-file`, `-o`: Write output to specified file (default: "-" for standard output).
  - `--directory`, `-D`: Add DIRECTORY to list for input files search.
  - `--update`, `-U`: Update def.po in-place; do nothing if def.po is already up to date.

- **Merging Options:**

  - `--compendium`, `-C`: Additional library of message translations (can be specified multiple times).
  - `--no-fuzzy-matching`, `-N`: Disable fuzzy matching (only use exact matches).
  - `--force-po`: Always write an output file even if empty.

- **Formatting Options:**

  - `--add-location`, `-n`: Generate '#: filename:line' lines (default: "full"). Options: `full`, `file`, or `never`.
  - `--no-location`: Suppress '#: filename:line' lines (same as `--add-location=never`).
  - `--no-wrap`: Do not break long message lines into multiple lines.
  - `--lang`: Set 'Language' field in the header entry (default: "en").

- **Help:**
  - `--help`, `-h`: Display help information.

### Examples

Basic merge with default options:

```bash
msgomerge translations.po messages.pot > merged.po
```

Merge and write directly to a file:

```bash
msgomerge -o merged.po translations.po messages.pot
```

Update the existing PO file in-place:

```bash
msgomerge -U translations.po messages.pot
```

Merge with strict matching (no fuzzy matching):

```bash
msgomerge -N -o merged.po translations.po messages.pot
```

Merge with additional translation libraries:

```bash
msgomerge -C common.po,extra.po -o merged.po translations.po messages.pot
```

Customize language and location tags:

```bash
msgomerge --lang=es --add-location=file -o merged.po translations.po messages.pot
```

## How It Works

1. **Input Files:**

   - `def.po`: An existing PO file with translations. These translations will be carried over to the new file if they still match the source strings.
   - `ref.pot`: A reference POT file (or PO file) with up-to-date source references but potentially outdated translations.

2. **Merging Process:**

   - Translations and comments from `def.po` are preserved.
   - Extracted comments and file positions from `def.po` are discarded.
   - Source references (file positions) from `ref.pot` are preserved.
   - Translations or comments in `ref.pot` are discarded.
   - Fuzzy matching is used when exact matches are not found (unless disabled).

3. **Output:**
   - The resulting PO file contains:
     - Preserved translations from `def.po` where they match the source strings.
     - Up-to-date source references from `ref.pot`.
     - Fuzzy matches (if enabled) for strings that don't have exact matches.

## Acknowledgments

- [gettext](https://www.gnu.org/software/gettext/) - The GNU internationalization and localization system that inspired this tool.
- [gotext](https://github.com/leonelquinteros/gotext) - The Go internationalization library this tool is designed to work with.
