package swag

import (
	"unicode"
)

type SwagParser struct {
	splitter *splitter
}

func (sp *SwagParser) Parse(line string) {
	_ = sp.splitter.split()

	return
}

type splitter struct {
	str           string
	maxSplitCount int
	pointer       int
}

func newSplitter(str string, maxSplitCount int) *splitter {
	return &splitter{str: str, maxSplitCount: maxSplitCount, pointer: -1}
}

func (s *splitter) split() []string {
	var result []string

	substr := []rune{}
	for {
		r, ok := s.next()
		if !ok {
			break
		}

		switch {
		case r == '"' && len(substr) == 0:
			result = append(result, s.splitQuoted())
		case unicode.IsSpace(r) || r == '\t':
			if len(substr) > 0 {
				result = append(result, string(substr))
				substr = []rune{}
			}
		default:
			substr = append(substr, r)
		}

		if s.maxSplitCount > 0 && len(result) == s.maxSplitCount-1 {
			result = append(result, string(s.str[s.pointer+1:]))
			break
		}
	}
	if len(substr) > 0 {
		result = append(result, string(substr))
	}

	return result
}

func (s *splitter) peek() (rune, bool) {
	if s.pointer >= len(s.str)-1 {
		return 0, false
	}

	return rune(s.str[s.pointer+1]), true
}

func (s *splitter) next() (rune, bool) {
	r, ok := s.peek()
	if !ok {
		return 0, false
	}

	s.pointer++
	return r, true
}

func (s *splitter) splitQuoted() string {
	substr := []rune{}
	for {
		r, ok := s.next()
		if !ok {
			break
		}

		if r == '"' {
			break
		}

		substr = append(substr, r)
	}

	return string(substr)
}
