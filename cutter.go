package goscut

import "unicode/utf8"

type Cutter struct {
	settings Settings
}

func NewCutter() *Cutter {
	c, err := NewCutterGivenSettings(nil)
	if err != nil {
		panic(err)
	}
	return c
}

func NewCutterGivenSettings(settings *Settings) (*Cutter, error) {
	if settings == nil {
		settings = &DefaultSettings
	}
	c := new(Cutter)
	settings.CopyTo(&c.settings)
	if c.settings.IsCutRune == nil {
		c.settings.IsCutRune = isSpaceFn
	}
	if err := c.settings.Check(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Cutter) Cut(s string) (
	result []string, surroundings [][2]rune, err error) {
	return c.cut(s, true, -1)
}

func (c *Cutter) CutKeepEmpty(s string) (
	result []string, surroundings [][2]rune, err error) {
	return c.cut(s, false, -1)
}

func (c *Cutter) CutN(s string, n int) (
	result []string, surroundings [][2]rune, err error) {
	return c.cut(s, true, n)
}

func (c *Cutter) CutNKeepEmpty(s string, n int) (
	result []string, surroundings [][2]rune, err error) {
	return c.cut(s, false, n)
}

func (c *Cutter) cut(s string, doesDiscardEmptyString bool, n int) (
	result []string, surroundings [][2]rune, err error) {
	if s == "" || n == 0 {
		return
	}
	if n == 1 {
		result = []string{s}
		surroundings = [][2]rune{{-1, -1}}
		return
	}
	if n > 0 {
		result = make([]string, 0, n)
		surroundings = make([][2]rune, 0, n)
	}
	var begin, end int
	var left, right rune = -1, -1
	var didAdd bool
	add := func() {
		didAdd = true
		if begin < 0 {
			return
		}
		if begin < end || (!doesDiscardEmptyString && begin == end) {
			result = append(result, s[begin:end])
			surroundings = append(surroundings, [2]rune{left, right})
		}
		left = right
		begin = -1 // To avoid adding the same content again after loop.
	}
	var badBegin, badEnd int = -1, -1
	null2Runes := [2]rune{-1, -1}
	quoteMarks, rawStringMarks := null2Runes, null2Runes
	var isEscaped, doesHit bool
	for i, r := range s {
		if didAdd {
			didAdd = false
			if len(result)+1 == n {
				// Add rest string.
				begin = end
				left, _ = utf8.DecodeLastRuneInString(s[:end])
				if left == utf8.RuneError {
					left = -1
				}
				end = len(s)
				right = -1
				add()
				return
			}
			begin = i
		}
		if err != nil {
			break
		}
		if r == utf8.RuneError {
			if badBegin < 0 {
				badBegin = i
			}
			continue
		} else if badBegin >= 0 {
			badEnd = i
			break
		}
		end = i
		right = r
		if rawStringMarks[1] != -1 {
			if r == rawStringMarks[1] {
				rawStringMarks = null2Runes
				add()
			}
		} else if quoteMarks[1] != -1 {
			if isEscaped {
				isEscaped = false
			} else if r == '\\' {
				isEscaped = true
			} else if r == quoteMarks[1] {
				quoteMarks = null2Runes
				add()
			}
		} else if c.settings.IsCutRune(r) {
			add()
		} else {
			for _, marks := range c.settings.QuoteMarks {
				if r == marks[0] {
					doesHit = true
					quoteMarks = marks
					add()
					break
				} else if r == marks[1] {
					doesHit = true
					err = NewHalfMarksError(marks[1], marks[0])
					add()
					break
				}
			}
			if doesHit {
				doesHit = false
				continue
			}
			for _, marks := range c.settings.RawStringMarks {
				if r == marks[0] {
					rawStringMarks = marks
					add()
					break
				} else if r == marks[1] {
					err = NewHalfMarksError(marks[1], marks[0])
					add()
					break
				}
			}
		}
	}
	if badBegin >= 0 {
		err = NewInvalidUtf8String(s, badBegin, badEnd)
	}
	end = len(s)
	right = -1
	add()
	if err == nil {
		if quoteMarks[0] != -1 {
			err = NewHalfMarksError(quoteMarks[0], quoteMarks[1])
		} else if rawStringMarks[0] != -1 {
			err = NewHalfMarksError(rawStringMarks[0], rawStringMarks[1])
		}
	}
	return
}
