package ast

import (
	"path/filepath"
	"slices"
	"strings"
)

type Entries []Entry

func (e Entries) ToNodes() []Node {
	var nodes []Node

	for _, entry := range e {
		nodes = append(nodes, entry)
	}

	return nodes
}

func (e Entries) AppendNodes(nodes ...Node) {
	for _, node := range nodes {
		if _, ok := node.(Entry); !ok {
			continue
		}
		e = append(e, node.(Entry))
	}
}

func (e Entries) Sort() Entries {
	groupsMap := make(map[string]Entries)
	for _, entry := range e {
		file := ""
		if len(entry.LocationComments) > 0 {
			file = filepath.Clean(entry.LocationComments[0].File)
		}
		groupsMap[file] = append(groupsMap[file], entry)
	}

	for k, group := range groupsMap {
		groupsMap[k] = group.SortByLine()
	}

	fileKeys := make([]string, 0, len(groupsMap))
	for file := range groupsMap {
		fileKeys = append(fileKeys, file)
	}
	slices.SortFunc(fileKeys, strings.Compare)

	var sortedEntries Entries
	for _, file := range fileKeys {
		sortedEntries = append(sortedEntries, groupsMap[file]...)
	}

	return sortedEntries
}

func (e Entries) SortByLine() Entries {
	sorted := make(Entries, len(e))
	copy(sorted, e)
	slices.SortFunc(sorted, func(a, b Entry) int {
		if len(a.LocationComments) == 0 {
			return 1
		}

		if len(b.LocationComments) == 0 {
			return -1
		}

		return a.LocationComments[0].Line - b.LocationComments[0].Line
	})

	return sorted
}

func (e Entries) CleanDuplicates() Entries {
	var cleaned Entries
	seedID := make(map[string][]int)

	for _, translation := range e {
		key := translation.Msgid.ID
		if translation.Msgctxt != nil {
			key += "\x00" + translation.Msgctxt.Context
		}

		if indices, exists := seedID[key]; exists {
			for _, idx := range indices {
				cleaned[idx].LocationComments = append(
					cleaned[idx].LocationComments,
					translation.LocationComments...)
			}
		} else {
			seedID[key] = []int{len(cleaned)}
			cleaned = append(cleaned, translation)
		}
	}

	return cleaned
}

func (e Entries) Solve() Entries {
	var cleaned Entries
	seedID := make(map[string][]int)

	for _, translation := range e {
		key := translation.Msgid.ID
		if translation.Msgctxt != nil {
			key += "\x00" + translation.Msgctxt.Context
		}

		if indices, exists := seedID[key]; exists {
			hasMsgstr := translation.Msgstr != nil && translation.Msgstr.Str != ""
			hasPlurals := len(translation.Plurals) > 0 && translation.Plurals[0].Str != ""

			for _, idx := range indices {
				if (hasMsgstr || hasPlurals) &&
					(cleaned[idx].Msgstr == nil || cleaned[idx].Msgstr.Str == "") &&
					len(cleaned[idx].Plurals) == 0 {
					cleaned[idx] = translation
				}
				cleaned[idx].LocationComments = append(
					cleaned[idx].LocationComments,
					translation.LocationComments...)
			}
		} else {
			seedID[key] = []int{len(cleaned)}
			cleaned = append(cleaned, translation)
		}
	}

	return cleaned
}

func MergeFiles(base *File, files ...*File) {
	for _, file := range files {
		base.Name += "_" + file.Name
		base.Nodes = append(base.Nodes, file.Nodes...)
	}

	var entries Entries
	entries.AppendNodes(base.Nodes...)

	base.Nodes = entries.Solve().Solve().ToNodes()
}
