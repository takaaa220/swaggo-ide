package swag

import (
	"iter"
	"unicode"

	"github.com/takaaa220/go-swag-ide/server/v2/server/util"
)

var (
	trimBraces = util.TrimBraces([][]rune{
		{'{', '}'},
		{'[', ']'},
		{'"', '"'},
	})
)

type splitter struct {
	str     string
	pointer int
}

func split(str string) (splitElement, func(maxSplitCount int) iter.Seq[splitElement]) {
	splitter := &splitter{str: str, pointer: -1}

	for {
		r, ok := splitter.peek()
		if !ok {
			return splitElement{}, nil
		}
		if !unicode.IsSpace(r) {
			break
		}
		splitter.next()
	}
	for {
		r, ok := splitter.peek()
		if !ok {
			return splitElement{}, nil
		}
		if r != '/' {
			break
		}
		splitter.next()
	}
	for {
		r, ok := splitter.peek()
		if !ok {
			return splitElement{}, nil
		}
		if !unicode.IsSpace(r) {
			break
		}
		splitter.next()
	}

	substr := []rune{}
	for {
		r, ok := splitter.peek()
		if !ok || unicode.IsSpace(r) {
			break
		}

		substr = append(substr, r)
		splitter.next()
	}

	return splitElement{
			Text:  string(substr),
			Start: splitter.pointer - len(substr) + 1,
			End:   splitter.pointer + 1,
		},
		func(maxSplitCount int) iter.Seq[splitElement] {
			return func(_yield func(splitElement) bool) {
				yield := func(res splitElement) bool {
					b := _yield(res)
					maxSplitCount--
					return b
				}

				substr := []rune{}
				if maxSplitCount == 1 {
					substr = splitter.getRest()
				}

				for {
					r, ok := splitter.peek()
					if !ok {
						break
					}

					yieldCalled := false
					switch {
					case r == '"' && len(substr) == 0:
						yield(splitter.splitSymbol('"', '"'))
						yieldCalled = true
					case unicode.IsSpace(r) || r == '\t':
						if len(substr) > 0 {
							yield(splitElement{
								Text:  string(substr),
								Start: splitter.pointer - len(substr) + 1,
								End:   splitter.pointer + 1,
							})
							yieldCalled = true
							substr = []rune{}
						}
						splitter.next()
					default:
						substr = append(substr, r)
						splitter.next()
					}

					if yieldCalled && maxSplitCount == 1 {
						substr = splitter.getRest()
					}
				}

				if len(substr) > 0 {
					yield(splitElement{
						Text:  string(substr),
						Start: splitter.pointer - len(substr) + 1,
						End:   splitter.pointer + 1,
					})
				}
			}
		}
}

type splitElement struct {
	Text  string
	Start int
	End   int
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

func (s *splitter) splitSymbol(openSymbol, closeSymbol rune) splitElement {
	r, ok := s.peek()
	if r != openSymbol || !ok {
		return splitElement{}
	}
	s.next()

	substr := []rune{r}
	for {
		r, ok := s.next()
		if !ok {
			break
		}

		substr = append(substr, r)
		if r == closeSymbol {
			break
		}
	}

	return splitElement{
		Text:  string(substr),
		Start: s.pointer - len(substr) + 1,
		End:   s.pointer + 1,
	}
}

func (s *splitter) getRest() []rune {
	substr := []rune{}
	for {
		r, ok := s.peek()
		if !ok {
			break
		}
		if !unicode.IsSpace(r) {
			break
		}
		s.next()
	}
	for {
		r, ok := s.peek()
		if !ok {
			break
		}
		s.next()
		substr = append(substr, r)
	}
	return substr
}
