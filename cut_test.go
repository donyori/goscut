package goscut

import "testing"

type TestCase struct {
	Input        string
	Result       []string
	Surroundings [][2]rune
	Error        error
}

var testCutCases = []TestCase{
	// Case 0:
	{
		"a b c",
		[]string{"a", "b", "c"},
		[][2]rune{
			{-1, ' '},
			{' ', ' '},
			{' ', -1},
		},
		nil,
	},
	// Case 1:
	{
		"a\tb \t c",
		[]string{"a", "b", "c"},
		[][2]rune{
			{-1, '\t'},
			{'\t', ' '},
			{' ', -1},
		},
		nil,
	},
	// Case 2:
	{
		"a \"b \t c\"",
		[]string{"a", "b \t c"},
		[][2]rune{
			{-1, ' '},
			{'"', '"'},
		},
		nil,
	},
	// Case 3:
	{
		"a \"b's c\"",
		[]string{"a", "b's c"},
		[][2]rune{
			{-1, ' '},
			{'"', '"'},
		},
		nil,
	},
	// Case 4:
	{
		"a \"\\b \\\"q\\\" c\"",
		[]string{"a", "\\b \\\"q\\\" c"},
		[][2]rune{
			{-1, ' '},
			{'"', '"'},
		},
		nil,
	},
	// Case 5:
	{
		"a `\"b's c\"`",
		[]string{"a", `"b's c"`},
		[][2]rune{
			{-1, ' '},
			{'`', '`'},
		},
		nil,
	},
	// Case 6:
	{
		"a `\\\"b's c\"`",
		[]string{"a", `\"b's c"`},
		[][2]rune{
			{-1, ' '},
			{'`', '`'},
		},
		nil,
	},
	// Case 7:
	{
		"a \"b c",
		[]string{"a", "b c"},
		[][2]rune{
			{-1, ' '},
			{'"', -1},
		},
		NewHalfMarksError('"', '"'),
	},
	// Case 8:
	{
		"a (b c",
		[]string{"a", "b c"},
		[][2]rune{
			{-1, ' '},
			{'(', -1},
		},
		NewHalfMarksError('(', ')'),
	},
	// Case 9:
	{
		"a )b c",
		[]string{"a", "b c"},
		[][2]rune{
			{-1, ' '},
			{')', -1},
		},
		NewHalfMarksError(')', '('),
	},
	// Case 10:
	{
		"a `b c",
		[]string{"a", "b c"},
		[][2]rune{
			{-1, ' '},
			{'`', -1},
		},
		NewHalfMarksError('`', '`'),
	},
	// Case 11:
	{
		"你好, 世界",
		[]string{"你好,", "世界"},
		[][2]rune{
			{-1, ' '},
			{' ', -1},
		},
		nil,
	},
	// Case 12:
	{
		"你好, 世界"[:13],
		[]string{"你好,", "世界"[:5]},
		[][2]rune{
			{-1, ' '},
			{' ', -1},
		},
		NewInvalidUtf8String("世界"[:5], 3, -1),
	},
	// Case 13:
	{
		"你好, 世界"[:13] + "后缀",
		[]string{"你好,", "世界"[:5] + "后缀"},
		[][2]rune{
			{-1, ' '},
			{' ', -1},
		},
		NewInvalidUtf8String("世界"[:5], 3, -1),
	},
}

var testCutKeepEmptyCases = []TestCase{
	{
		"a\tb \t c",
		[]string{"a", "b", "", "", "c"},
		[][2]rune{
			{-1, '\t'},
			{'\t', ' '},
			{' ', '\t'},
			{'\t', ' '},
			{' ', -1},
		},
		nil,
	},
	{
		"a\tb \"\"\t c",
		[]string{"a", "b", "", "", "", "", "c"},
		[][2]rune{
			{-1, '\t'},
			{'\t', ' '},
			{' ', '"'},
			{'"', '"'},
			{'"', '\t'},
			{'\t', ' '},
			{' ', -1},
		},
		nil,
	},
}

var testCutNCases = []TestCase{
	{
		"a   b c",
		[]string{"a", "b c"},
		[][2]rune{
			{-1, ' '},
			{' ', -1},
		},
		nil,
	},
}

var testCutNKeepEmptyCases = []TestCase{
	{
		"a   b c",
		[]string{"a", "  b c"},
		[][2]rune{
			{-1, ' '},
			{' ', -1},
		},
		nil,
	},
	{
		"a\tb \"\"\t c",
		[]string{"a", "b", "", "\"\"\t c"},
		[][2]rune{
			{-1, '\t'},
			{'\t', ' '},
			{' ', '"'},
			{' ', -1},
		},
		nil,
	},
}

func TestCut(t *testing.T) {
	testCore(t, testCutCases, func(i int, s string) ([]string, [][2]rune, error) {
		return Cut(s)
	})
}

func TestCutKeepEmpty(t *testing.T) {
	testCore(t, testCutKeepEmptyCases, func(i int, s string) ([]string, [][2]rune, error) {
		return CutKeepEmpty(s)
	})
}

func TestCutN(t *testing.T) {
	testCore(t, testCutNCases, func(i int, s string) ([]string, [][2]rune, error) {
		return CutN(s, 2)
	})
}

func TestCutNKeepEmpty(t *testing.T) {
	testCore(t, testCutNKeepEmptyCases, func(i int, s string) ([]string, [][2]rune, error) {
		if i == 0 {
			return CutNKeepEmpty(s, 2)
		}
		return CutNKeepEmpty(s, 4)
	})
}

func testCore(t *testing.T, cases []TestCase, fn func(int, string) ([]string, [][2]rune, error)) {
	for i, c := range cases {
		t.Log("Case", i)
		r, s, e := fn(i, c.Input)
		t.Log(c.Input, r, s, e)
		if len(r) != len(c.Result) || (r == nil) != (c.Result == nil) {
			t.Error("Wrong result, want", c.Result, "got", r)
		} else {
			for i := range r {
				if r[i] != c.Result[i] {
					t.Error("Wrong result, want", c.Result, "got", r)
				}
			}
		}
		if len(s) != len(c.Surroundings) || (s == nil) != (c.Surroundings == nil) {
			t.Error("Wrong surroundings, want", c.Surroundings, "got", s)
		} else {
			for i := range s {
				if s[i] != c.Surroundings[i] {
					t.Error("Wrong surroundings, want", c.Surroundings, "got", s)
				}
			}
		}
		if (e == nil) != (c.Error == nil) {
			t.Error("Unexpected error, want", c.Error, "got", e)
		} else if e != nil && e.Error() != c.Error.Error() {
			t.Error("Unexpected error, want", c.Error, "got", e)
		}
	}
}
