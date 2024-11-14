package swag

import (
	"iter"
	"unicode"
)

type tokenizer struct {
	str     string
	pointer int
}

func tokenize(str string) (token, func(maxTokenCount int) iter.Seq[token]) {
	tokenizer := &tokenizer{str: str, pointer: -1}

	for {
		r, ok := tokenizer.peek()
		if !ok {
			return token{}, nil
		}
		if !unicode.IsSpace(r) {
			break
		}
		tokenizer.next()
	}
	for {
		r, ok := tokenizer.peek()
		if !ok {
			return token{}, nil
		}
		if r != '/' {
			break
		}
		tokenizer.next()
	}
	for {
		r, ok := tokenizer.peek()
		if !ok {
			return token{}, nil
		}
		if !unicode.IsSpace(r) {
			break
		}
		tokenizer.next()
	}

	substr := []rune{}
	for {
		r, ok := tokenizer.peek()
		if !ok || unicode.IsSpace(r) {
			break
		}

		substr = append(substr, r)
		tokenizer.next()
	}

	return token{
			Text:  string(substr),
			Start: tokenizer.pointer - len(substr) + 1,
			End:   tokenizer.pointer + 1,
		},
		func(maxTokenCount int) iter.Seq[token] {
			return func(_yield func(token) bool) {
				yield := func(res token) bool {
					b := _yield(res)
					maxTokenCount--
					return b
				}

				substr := []rune{}
				if maxTokenCount == 1 {
					substr = tokenizer.getRest()
				}

				for {
					r, ok := tokenizer.peek()
					if !ok {
						break
					}

					yieldCalled := false
					switch {
					case r == '"' && len(substr) == 0:
						if !yield(tokenizer.splitSymbol('"', '"')) {
							return
						}
						yieldCalled = true
					case unicode.IsSpace(r) || r == '\t':
						if len(substr) > 0 {
							if !yield(token{
								Text:  string(substr),
								Start: tokenizer.pointer - len(substr) + 1,
								End:   tokenizer.pointer + 1,
							}) {
								return
							}
							yieldCalled = true
							substr = []rune{}
						}
						tokenizer.next()
					default:
						substr = append(substr, r)
						tokenizer.next()
					}

					if yieldCalled && maxTokenCount == 1 {
						substr = tokenizer.getRest()
					}
				}

				if len(substr) > 0 {
					if !yield(token{
						Text:  string(substr),
						Start: tokenizer.pointer - len(substr) + 1,
						End:   tokenizer.pointer + 1,
					}) {
						return
					}
				}
			}
		}
}

type token struct {
	Text  string
	Start int
	End   int
}

func (s *tokenizer) peek() (rune, bool) {
	if s.pointer >= len(s.str)-1 {
		return 0, false
	}

	return rune(s.str[s.pointer+1]), true
}

func (s *tokenizer) next() (rune, bool) {
	r, ok := s.peek()
	if !ok {
		return 0, false
	}

	s.pointer++
	return r, true
}

func (s *tokenizer) splitSymbol(openSymbol, closeSymbol rune) token {
	r, ok := s.peek()
	if r != openSymbol || !ok {
		return token{}
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

	return token{
		Text:  string(substr),
		Start: s.pointer - len(substr) + 1,
		End:   s.pointer + 1,
	}
}

func (s *tokenizer) getRest() []rune {
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
