package goscut

func Cut(s string) (result []string, surroundings [][2]rune, err error) {
	c, err := NewCutterGivenSettings(nil)
	if err != nil {
		return nil, nil, err
	}
	return c.Cut(s)
}

func CutKeepEmpty(s string) (
	result []string, surroundings [][2]rune, err error) {
	c, err := NewCutterGivenSettings(nil)
	if err != nil {
		return nil, nil, err
	}
	return c.CutKeepEmpty(s)
}

func CutN(s string, n int) (
	result []string, surroundings [][2]rune, err error) {
	c, err := NewCutterGivenSettings(nil)
	if err != nil {
		return nil, nil, err
	}
	return c.CutN(s, n)
}

func CutNKeepEmpty(s string, n int) (
	result []string, surroundings [][2]rune, err error) {
	c, err := NewCutterGivenSettings(nil)
	if err != nil {
		return nil, nil, err
	}
	return c.CutNKeepEmpty(s, n)
}
