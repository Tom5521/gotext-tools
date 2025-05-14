package po

import (
	"errors"
	"fmt"
)

var ErrBadPluralEntry = errors.New("the entry can't be plural and singular at the same time")

type DuplicatedEntryError struct {
	OriginalIndex int
}

func (e *DuplicatedEntryError) Error() string {
	return fmt.Sprintf("entry is a duplicated of %d", e.OriginalIndex)
}

type InvalidEntryError struct {
	ID     string
	Reason error
}

func (e *InvalidEntryError) Error() string {
	return fmt.Sprintf("entry %q is invalid: %v", e.ID, e.Reason)
}

func (e *InvalidEntryError) Unwrap() error {
	return e.Reason
}

type InvalidEntryAtIndexError struct {
	Index  int
	Reason error
}

func (e *InvalidEntryAtIndexError) Error() string {
	return fmt.Sprintf("invalid entry at %d: %v", e.Index, e.Reason)
}

func (e *InvalidEntryAtIndexError) Unwrap() error {
	return e.Reason
}

type InvalidFileError struct {
	Filename string
	Reason   error
}

func (e *InvalidFileError) Error() string {
	return fmt.Sprintf("file %q is invalid: %v", e.Filename, e.Reason)
}

func (e *InvalidFileError) Unwrap() error {
	return e.Reason
}
