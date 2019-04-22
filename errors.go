package goscut

import (
	"errors"
	"fmt"
)

type RecurringMarksError struct {
	marks [2]rune
}

type MarkInCutSetError struct {
	mark rune
}

type InvalidUtf8String struct {
	b []byte
}

type HalfMarksError struct {
	got, missed rune
}

var (
	ErrBackslashInCutSet error = errors.New(`goscut: cannot use "\" in cut set`)
	ErrBackslashAsMark   error = errors.New(`goscut: cannot use "\" as a mark`)
)

func NewRecurringMarksError(marks [2]rune) error {
	return &RecurringMarksError{marks: marks}
}

func (rme *RecurringMarksError) Error() string {
	return fmt.Sprintf("goscut: found recurring mark %c%c",
		rme.marks[0], rme.marks[1])
}

func NewMarkInCutSetError(mark rune) error {
	return &MarkInCutSetError{mark: mark}
}

func (micse *MarkInCutSetError) Error() string {
	return fmt.Sprintf("goscut: mark %c is in cut set", micse.mark)
}

func NewInvalidUtf8String(s string, start, end int) error {
	if start < 0 {
		start = 0
	}
	if end < 0 || end > len(s) {
		end = len(s)
	}
	return &InvalidUtf8String{b: []byte(s[start:end])}
}

func (ius *InvalidUtf8String) Error() string {
	return fmt.Sprintf("goscut: UTF-8 string is invalid %q", ius.b)
}

func NewHalfMarksError(got, missed rune) error {
	return &HalfMarksError{got: got, missed: missed}
}

func (hme *HalfMarksError) Error() string {
	return fmt.Sprintf("goscut: half of marks missed, got %c, missed %c",
		hme.got, hme.missed)
}
