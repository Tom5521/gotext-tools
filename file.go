package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// MsgID represents a message identifier with its line number.
type MsgID struct {
	File string
	ID   string
	Line int
}

// File handles the processing of Go source files for message extraction.
type File struct {
	Path       string
	HasImport  bool
	ImportName string
	rawData    []byte
	content    string
	imports    string
	getsRegex  *regexp.Regexp
}

// NewFile creates and initializes a new File instance.
func NewFile(path string) (*File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %w", path, err)
	}

	f := &File{
		Path:    path,
		rawData: data,
		content: string(data),
	}

	f.parseImports()

	if f.HasImport {
		f.getsRegex = regexp.MustCompile(fmt.Sprintf(baseRegex, f.ImportName))
	}

	return f, nil
}

func (f *File) parseImports() {
	importedPkgs := importsCompiler.FindAllString(f.content, -1)
	f.imports = f.buildImportsString(importedPkgs)
	f.HasImport = f.checkForWantedPackage(importedPkgs)

	if f.HasImport {
		f.ImportName = f.extractImportName()
	}
}

func (f *File) buildImportsString(pkgs []string) string {
	var builder strings.Builder
	for _, pkg := range pkgs {
		builder.WriteString(pkg)
	}
	return builder.String()
}

func (f *File) checkForWantedPackage(pkgs []string) bool {
	for _, match := range pkgs {
		if strings.Contains(match, wantedPkg) {
			return true
		}
	}
	return false
}

func (f *File) extractImportName() string {
	for _, line := range strings.Split(f.imports, "\n") {
		if !strings.Contains(line, wantedPkg) {
			continue
		}

		parts := strings.SplitN(line, `"`, 2)
		if len(parts) > 0 && !isEmpty(parts[0]) {
			return cleanWhitespaces(parts[0])
		}
		return "gotext"
	}
	return "gotext"
}

// MsgIDs extracts all message IDs from the file.
func (f *File) MsgIDs() []MsgID {
	if !f.HasImport || f.getsRegex == nil {
		return nil
	}

	var msgids []MsgID
	indexes := f.getsRegex.FindAllIndex(f.rawData, -1)

	for _, index := range indexes {
		msgids = append(msgids, MsgID{
			Line: findLine(f.content, index[0]),
			ID:   strContent(f.content[index[0]:index[1]]),
			File: f.Path,
		})
	}

	return msgids
}
