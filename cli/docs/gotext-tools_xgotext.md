## gotext-tools xgotext

Extract translatable strings from given input files.

### Synopsis

Extract translatable strings from given input files.
Mandatory arguments to long options are mandatory for short options too.
Similarly for optional arguments.

```
gotext-tools xgotext [flags]
```

### Examples

```
xgotext -a my-go-proyect/ -o -
xgotext -a my-go-proyect/ -o en.pot -lang en
xgotext main.go -o main.pot -lang en
```

### Options

```
      --add-location string         Generate ‘#: filename:line’ lines (default).
                                    
                                    The optional type can be either ‘full’, ‘file’, or ‘never’. 
                                    If it is not given or ‘full’, it generates the lines with both file
                                    name and line number. If it is ‘file’, the line number part is omitted. 
                                    If it is ‘never’, it completely suppresses the lines (same as --no-location).
                                     (default "full")
      --copyright-holder string     Set the copyright holder in the output.
                                    string should be the copyright holder of the surrounding package.
                                    (Note that the msgstr strings, extracted from the package’s sources,
                                    belong to the copyright holder of the package.) Translators are expected to transfer
                                    or disclaim the copyright for their translations, so that package maintainers
                                    can distribute them without legal risk. If string is empty, the output files
                                    are marked as being in the public domain; in this case, the translators are
                                    expected to disclaim their copyright, again so that package maintainers can
                                    distribute them without legal risk.
                                     (default "YEAR THE PACKAGE'S COPYRIGHT HOLDER")
  -d, --default-domain string       use NAME.pot for output (instead of messages.pot) (default "messages")
  -D, --directory string            add DIRECTORY to list for input files search
                                    If input file is -, standard input is read.
  -X, --exclude strings             Specifies which files will be omitted.
  -x, --exclude-file string         Entries from file are not extracted. file should be a PO or POT file.
  -a, --extract-all                 Extract all strings.
  -f, --files-from string           get list of input files from FILE
      --force-po                    Always write an output file even if no message is defined.
      --foreign-user                Omit FSF copyright in output. 
                                    This option is equivalent to ‘--copyright-holder=''’. It can be useful
                                    for packages outside the GNU project that want their translations to
                                    be in the public domain.
  -h, --help                        help for xgotext
  -j, --join-existing               Join messages with existing file.
  -l, --lang string                 Language code to include in the POT file (default "en")
      --msgid-bugs-address string   Set the reporting address for msgid bugs.
                                    This is the email address or URL to which the translators shall report bugs
                                    in the untranslated strings:
                                        - Strings which are not entire sentences; see the maintainer guidelines
                                        in Preparing Translatable Strings.
                                        - Strings which use unclear terms or require additional context to be understood.
                                        - Strings which make invalid assumptions about notation of date, time or money.
                                        - Pluralisation problems.
                                        - Incorrect English spelling.
                                        - Incorrect formatting. 
                                    
                                    It can be your email address, or a mailing list address where translators 
                                    can write to without being subscribed, or the URL of a web page through which
                                    the translators can contact you.
                                    
                                    The default value is empty, which means that translators will be clueless!
                                    Don’t forget to specify this option.
                                    
  -m, --msgstr-prefix string        Use string (or "" if not specified) as prefix for msgstr values.
  -M, --msgstr-suffix string        Use string (or "" if not specified) as suffix for msgstr values.
  -n, --no-location                 Do not write ‘#: filename:line’ lines.
                                    Note that using this option makes it harder for technically
                                    skilled translators to understand each message’s context.
                                    
      --nplurals uint               Specify the number of plurals forms of the language in question (default 2)
      --omit-header                 Don’t write header with ‘msgid ""’ entry. 
                                    Note: Using this option may lead to an error in subsequent
                                    operations if the output contains non-ASCII characters.
                                    
  -o, --output string               write output to specified file (default "-")
  -p, --output-dir string           output files will be placed in directory DIR
                                    If output file is -, output is written to standard output.
      --package-name string         Set the package name in the header of the output. (default "PACKAGE")
      --package-version string      Set the package version in the header of the output. 
                                    This option has an effect only if the ‘--package-name’ option is also used. (default "PACKAGE VERSION")
      --title string                Set the title of the pot file. (default "SOME DESCRIPTIVE TITLE")
      --verbose                     increase verbosity level
      --word-wrap                   Applies word wrapping to strings.
```

### SEE ALSO

* [gotext-tools](gotext-tools.md)	 - A wrapper for the CLI tools from github.com/Tom5521/gotext-tools/v2/cli

###### Auto generated by spf13/cobra on 5-May-2025
