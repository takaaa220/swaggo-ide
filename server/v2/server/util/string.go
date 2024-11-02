package util

func TrimBraces(symbolPairs [][]rune) func(s string) string {
	return func(s string) string {
		if len(s) < 2 {
			return s
		}

		for _, symbolPair := range symbolPairs {
			if len(symbolPair) != 2 {
				continue
			}
			if rune(s[0]) == symbolPair[0] && rune(s[len(s)-1]) == symbolPair[1] {
				return s[1 : len(s)-1]
			}
		}

		return s
	}
}
