package goscut

import "unicode"

type IsCutRuneFn func(r rune) bool

type Settings struct {
	IsCutRune      IsCutRuneFn
	QuoteMarks     [][2]rune
	RawStringMarks [][2]rune
}

var DefaultSettings Settings

var isSpaceFn IsCutRuneFn = unicode.IsSpace

func init() {
	DefaultSettings.Reset()
}

func IsCutRuneFromString(cutSet string) IsCutRuneFn {
	m := make(map[rune]bool)
	for _, r := range cutSet {
		m[r] = true
	}
	f := IsCutRuneFn(func(r rune) bool {
		return m[r]
	})
	return f
}

func NewSettings() *Settings {
	s := new(Settings)
	s.Reset()
	return s
}

func (s *Settings) Reset() {
	s.IsCutRune = isSpaceFn
	s.QuoteMarks = [][2]rune{
		{'\'', '\''},
		{'"', '"'},
		{'(', ')'},
		{'[', ']'},
		{'{', '}'},
		{'<', '>'},
	}
	s.RawStringMarks = [][2]rune{
		{'`', '`'},
	}
}

func (s *Settings) Copy() *Settings {
	if s == nil {
		return nil
	}
	c := new(Settings)
	s.CopyTo(c)
	return c
}

func (s *Settings) CopyTo(settings *Settings) {
	if s == nil {
		if settings == nil {
			return
		}
		settings.IsCutRune = nil
		settings.QuoteMarks = nil
		settings.RawStringMarks = nil
		return
	}
	settings.IsCutRune = s.IsCutRune
	// Copy slices. See https://github.com/go101/go101/wiki for details.
	settings.QuoteMarks = append(s.QuoteMarks[:0:0], s.QuoteMarks...)
	settings.RawStringMarks = append(s.RawStringMarks[:0:0], s.RawStringMarks...)
}

func (s *Settings) Check() error {
	if s == nil {
		return nil
	}
	if s.IsCutRune != nil && s.IsCutRune('\\') {
		return ErrBackslashInCutSet
	}
	occurred := make(map[[2]rune]bool)
	for _, marks := range s.QuoteMarks {
		if marks[0] == '\\' || marks[1] == '\\' {
			return ErrBackslashAsMark
		}
		if s.IsCutRune != nil {
			if s.IsCutRune(marks[0]) {
				return NewMarkInCutSetError(marks[0])
			}
			if s.IsCutRune(marks[1]) {
				return NewMarkInCutSetError(marks[1])
			}
		}
		if occurred[marks] {
			return NewRecurringMarksError(marks)
		}
		occurred[marks] = true
	}
	for _, marks := range s.RawStringMarks {
		if marks[0] == '\\' || marks[1] == '\\' {
			return ErrBackslashAsMark
		}
		if s.IsCutRune != nil {
			if s.IsCutRune(marks[0]) {
				return NewMarkInCutSetError(marks[0])
			}
			if s.IsCutRune(marks[1]) {
				return NewMarkInCutSetError(marks[1])
			}
		}
		if occurred[marks] {
			return NewRecurringMarksError(marks)
		}
		occurred[marks] = true
	}
	return nil
}
