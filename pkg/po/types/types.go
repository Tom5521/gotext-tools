package types

type HeaderField struct {
	Key   string
	Value string
}

type Header struct {
	Values []HeaderField
}

type File struct {
	Name         string
	Header       Header
	Translations []Translation
}

// Location represents the location of a translation string in the source code.
type Location struct {
	Line int    // The line number of the translation.
	File string // The file name where the translation is located.
}

// Translation represents a translatable string, including its context, plural forms,
// and source code locations.
type Translation struct {
	ID        string     // The original string to be translated.
	Context   string     // The context in which the string is used (optional).
	Plural    string     // The plural form of the string (optional).
	Locations []Location // A list of source code locations for the string.
}
