package swag

import (
	"fmt"
	"strings"
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

type SwagChecker struct {
}

func NewSwagChecker() *SwagChecker {
	return &SwagChecker{}
}

func (sp *SwagChecker) Check(line string) (bool, []string) {
	line = strings.TrimSpace(line)

	splitForTag := newSplitter(line, 2).split()
	if len(splitForTag) == 0 {
		return false, []string{}
	}

	tag := splitForTag[0]
	swagTagDef := newSwagTagDef(strings.TrimPrefix(tag, "@"))
	if swagTagDef._type == swagTagTypeUnknown {
		return false, []string{}
	}

	splitArgs := []string{}
	if len(splitForTag) > 1 {
		splitArgs = newSplitter(splitForTag[1], len(swagTagDef.args)).split()
	}

	if len(splitArgs) < swagTagDef.requiredArgsCount {
		return false, []string{swagTagDef.errorMessage()}
	}

	errors := []string{}

	for i, argDef := range swagTagDef.args {
		if i >= len(splitArgs) {
			break
		}

		argStr := trimBraces(splitArgs[i])

		var arg swagTagArg

		switch argDef.valueType {
		case swagTagArgDefTypeString:
			arg = &swagTagArgString{value: argStr}
		case swagTagArgDefTypeGoType:
			arg = &swagTagArgGoType{value: argStr}
		default:
			panic(fmt.Errorf("unknown argDef.valueType: %d", argDef.valueType))
		}

		ok, errorMessages := argDef.isValid(arg)
		if !ok {
			errors = append(errors, fmt.Sprintf("%s(Arg=%d): %s", argDef.label, i+1, strings.Join(errorMessages, ", ")))
		}
	}

	return len(errors) == 0, errors
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
	if s.maxSplitCount == 1 {
		return []string{s.str}
	}

	var result []string

	substr := []rune{}
	for {
		r, ok := s.peek()
		if !ok {
			break
		}

		switch {
		case r == '"' && len(substr) == 0:
			result = append(result, s.splitSymbol('"', '"'))
		case r == '{' && len(substr) == 0:
			result = append(result, s.splitSymbol('{', '}'))
		case r == '[' && len(substr) == 0:
			result = append(result, s.splitSymbol('[', ']'))
		case unicode.IsSpace(r) || r == '\t':
			if len(substr) > 0 {
				result = append(result, string(substr))
				substr = []rune{}
			}
			s.next()
		default:
			substr = append(substr, r)
			s.next()
		}

		if s.maxSplitCount > 0 && len(result) == s.maxSplitCount-1 {
			result = append(result, strings.TrimSpace(s.str[s.pointer+1:]))
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

func (s *splitter) splitSymbol(openSymbol, closeSymbol rune) string {
	r, ok := s.peek()
	if r != openSymbol || !ok {
		return ""
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

	return string(substr)
}
