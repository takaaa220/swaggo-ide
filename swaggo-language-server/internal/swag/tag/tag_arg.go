package tag

import (
	"fmt"
	"strconv"
	"strings"
)

type swagTagArgDefType int

const (
	swagTagArgDefTypeString swagTagArgDefType = iota
	swagTagArgDefTypeGoType
)

type SwagTagArgDef struct {
	valueType swagTagArgDefType
	label     string
	checkers  []swagTagArgChecker
}

func newSwagTagStringArgDef(label string, optionalCheckers ...swagTagArgChecker) SwagTagArgDef {
	return SwagTagArgDef{
		valueType: swagTagArgDefTypeString,
		label:     label,
		checkers:  append([]swagTagArgChecker{&swagTagArgStringChecker{}}, optionalCheckers...),
	}
}

func newSwagTagIntArgDef(label string, optionalCheckers ...swagTagArgChecker) SwagTagArgDef {
	return SwagTagArgDef{
		valueType: swagTagArgDefTypeString,
		label:     label,
		checkers:  append([]swagTagArgChecker{&swagTagArgIntChecker{}}, optionalCheckers...),
	}
}

func newSwagTagBoolArgDef(label string, optionalCheckers ...swagTagArgChecker) SwagTagArgDef {
	return SwagTagArgDef{
		valueType: swagTagArgDefTypeString,
		label:     label,
		checkers:  append([]swagTagArgChecker{&swagTagArgBoolChecker}, optionalCheckers...),
	}
}

func newSwagTagGoTypeArgDef(label string, optionalCheckers ...swagTagArgChecker) SwagTagArgDef {
	return SwagTagArgDef{
		valueType: swagTagArgDefTypeGoType,
		label:     label,
		checkers:  append([]swagTagArgChecker{&swagTagArgGoTypeChecker{}}, optionalCheckers...),
	}
}

func newSwagTagConstStringArgDef(label string, value string, optionalCheckers ...swagTagArgChecker) SwagTagArgDef {
	return SwagTagArgDef{
		valueType: swagTagArgDefTypeString,
		label:     label,
		checkers:  append([]swagTagArgChecker{&swagTagArgConstStringChecker{value: value}}, optionalCheckers...),
	}
}

func newSwagTagUnionArgDef(label string, checker *swagTagArgUnionChecker, optionalCheckers ...swagTagArgChecker) SwagTagArgDef {
	return SwagTagArgDef{
		valueType: swagTagArgDefTypeString,
		label:     label,
		checkers:  append([]swagTagArgChecker{checker}, optionalCheckers...),
	}
}

func wrapArgDefWithBraces(openBrace rune, argDef SwagTagArgDef) SwagTagArgDef {
	braces := map[rune]rune{
		'[': ']',
		'{': '}',
		'"': '"',
	}
	closeBrace, ok := braces[openBrace]
	if !ok {
		panic(fmt.Errorf("invalid openBrace: %c", openBrace))
	}

	checkers := argDef.checkers
	for i, checker := range checkers {

		checkers[i] = &swagTagArgWithBracesChecker{
			openBrace:  openBrace,
			closeBrace: closeBrace,
			checker:    checker,
		}
	}

	return SwagTagArgDef{
		valueType: argDef.valueType,
		label:     string(openBrace) + argDef.Label() + string(closeBrace),
		checkers:  checkers,
	}
}

func (s *SwagTagArgDef) Check(arg swagTagArg) (bool, []string) {
	errors := []string{}
	for _, checker := range s.checkers {
		if ok, errorMessage := checker.check(arg); !ok {
			errors = append(errors, fmt.Sprintf("%s %s.", s.Label(), errorMessage))
		}
	}

	return len(errors) == 0, errors
}

func (s *SwagTagArgDef) Candidates() []string {
	candidates := []string{}
	for _, checker := range s.checkers {
		candidates = append(candidates, checker.candidates()...)
	}
	return candidates
}

func (s *SwagTagArgDef) Label() string {
	return s.label
}

type swagTagArgChecker interface {
	// check returns check result and error message
	check(arg swagTagArg) (bool, string)
	candidates() []string
}

var (
	_ swagTagArgChecker = &swagTagArgStringChecker{}
	_ swagTagArgChecker = &swagTagArgIntChecker{}
	_ swagTagArgChecker = &swagTagArgConstStringChecker{}
	_ swagTagArgChecker = &swagTagArgUnionChecker{}
	_ swagTagArgChecker = &swagTagArgGoTypeChecker{}
	_ swagTagArgChecker = &swagTagArgWithBracesChecker{}
)

type swagTagArgStringChecker struct{}

func (s *swagTagArgStringChecker) check(arg swagTagArg) (bool, string) {
	if _, ok := arg.(*swagTagArgString); !ok {
		return false, "invalid string"
	}

	return true, ""
}

func (s *swagTagArgStringChecker) candidates() []string {
	return []string{}
}

type swagTagArgWithBracesChecker struct {
	openBrace  rune
	closeBrace rune
	checker    swagTagArgChecker
}

func (s *swagTagArgWithBracesChecker) check(arg swagTagArg) (bool, string) {
	stringArg, ok := arg.(*swagTagArgString)
	if !ok {
		return false, "should be string"
	}

	if len(stringArg.value) < 2 || rune(stringArg.value[0]) != s.openBrace || rune(stringArg.value[len(stringArg.value)-1]) != s.closeBrace {
		return false, "should be enclosed in braces"
	}

	return s.checker.check(&swagTagArgString{value: stringArg.value[1 : len(stringArg.value)-1]})
}

func (s *swagTagArgWithBracesChecker) candidates() []string {
	candidates := s.checker.candidates()
	for i, candidate := range candidates {
		candidates[i] = string(s.openBrace) + candidate + string(s.closeBrace)
	}
	return candidates
}

type swagTagArgIntChecker struct{}

func (s *swagTagArgIntChecker) check(arg swagTagArg) (bool, string) {
	v, ok := arg.(*swagTagArgString)
	if !ok {
		return false, "should be integer"
	}

	_, err := strconv.Atoi(v.value)
	if err != nil {
		return false, "should be integer"
	}

	return true, ""
}

func (s *swagTagArgIntChecker) candidates() []string {
	return []string{}
}

type swagTagArgConstStringChecker struct {
	value     string
	converter func(string) string
}

func (s *swagTagArgConstStringChecker) check(arg swagTagArg) (bool, string) {
	stringArg, ok := arg.(*swagTagArgString)
	if !ok {
		return false, "should be string"
	}

	argStr := stringArg.value
	if s.converter != nil {
		argStr = s.converter(argStr)
	}

	if argStr != s.value {
		return false, "should be " + s.value
	}
	return true, ""
}

func (s *swagTagArgConstStringChecker) candidates() []string {
	return []string{s.value}
}

type swagTagArgUnionChecker struct {
	options      []swagTagArgChecker
	errorMessage string
}

func (s *swagTagArgUnionChecker) check(arg swagTagArg) (bool, string) {
	for _, option := range s.options {
		if ok, _ := option.check(arg); ok {
			return true, ""
		}
	}

	return false, s.errorMessage
}

func (s *swagTagArgUnionChecker) candidates() []string {
	candidates := []string{}
	for _, option := range s.options {
		candidates = append(candidates, option.candidates()...)
	}
	return candidates
}

type swagTagArgGoTypeChecker struct{}

func (s *swagTagArgGoTypeChecker) check(arg swagTagArg) (bool, string) {
	// TODO: Implement
	return true, ""
}

func (s *swagTagArgGoTypeChecker) candidates() []string {
	return []string{}
}

var (
	userDefinedType                = &swagTagArgGoTypeChecker{}
	swagTagArgMimeTypeUnionChecker = &swagTagArgUnionChecker{
		options: []swagTagArgChecker{
			&swagTagArgConstStringChecker{value: "json"},
			&swagTagArgConstStringChecker{value: "application/json"},
			&swagTagArgConstStringChecker{value: "xml"},
			&swagTagArgConstStringChecker{value: "text/xml"},
			&swagTagArgConstStringChecker{value: "plain"},
			&swagTagArgConstStringChecker{value: "text/plain"},
			&swagTagArgConstStringChecker{value: "html"},
			&swagTagArgConstStringChecker{value: "text/html"},
			&swagTagArgConstStringChecker{value: "mpfd"},
			&swagTagArgConstStringChecker{value: "multipart/form-data"},
			&swagTagArgConstStringChecker{value: "x-www-form-urlencoded"},
			&swagTagArgConstStringChecker{value: "application/x-www-form-urlencoded"},
			&swagTagArgConstStringChecker{value: "json-api"},
			&swagTagArgConstStringChecker{value: "application/vnd.api+json"},
			&swagTagArgConstStringChecker{value: "json-stream"},
			&swagTagArgConstStringChecker{value: "application/x-json-stream"},
			&swagTagArgConstStringChecker{value: "octet-stream"},
			&swagTagArgConstStringChecker{value: "application/octet-stream"},
			&swagTagArgConstStringChecker{value: "png"},
			&swagTagArgConstStringChecker{value: "image/png"},
			&swagTagArgConstStringChecker{value: "jpeg"},
			&swagTagArgConstStringChecker{value: "image/jpeg"},
			&swagTagArgConstStringChecker{value: "gif"},
			&swagTagArgConstStringChecker{value: "image/gif"},
		},
		errorMessage: "should be valid mime type",
	}
	swagTagArgParamTypeUnionChecker = &swagTagArgUnionChecker{
		options: []swagTagArgChecker{
			&swagTagArgConstStringChecker{value: "path"},
			&swagTagArgConstStringChecker{value: "query"},
			&swagTagArgConstStringChecker{value: "header"},
			&swagTagArgConstStringChecker{value: "body"},
			&swagTagArgConstStringChecker{value: "formData"},
			&swagTagArgConstStringChecker{value: "object"},
		},
		errorMessage: "should be `path, query, header, body, formData, or object`",
	}
	swagTagArgDataTypeUnionChecker = &swagTagArgUnionChecker{
		options: []swagTagArgChecker{
			&swagTagArgConstStringChecker{value: "string"},
			&swagTagArgConstStringChecker{value: "number"},
			&swagTagArgConstStringChecker{value: "integer"},
			&swagTagArgConstStringChecker{value: "boolean"},
			&swagTagArgConstStringChecker{value: "file"},
			&swagTagArgConstStringChecker{value: "object"},
		},
		errorMessage: "should be `string, number, integer, boolean, file or object`",
	}
	swagTagArgGoDataTypeUnionChecker = &swagTagArgUnionChecker{
		options: []swagTagArgChecker{
			&swagTagArgConstStringChecker{value: "string"},
			&swagTagArgConstStringChecker{value: "number"},
			&swagTagArgConstStringChecker{value: "integer"},
			&swagTagArgConstStringChecker{value: "boolean"},
			&swagTagArgConstStringChecker{value: "file"},
			&swagTagArgConstStringChecker{value: "object"},
			userDefinedType,
		},
		errorMessage: "should be `string, number, integer, boolean, file, object or user-defined type`",
	}
	swagTagArgBoolChecker = swagTagArgUnionChecker{
		options: []swagTagArgChecker{
			&swagTagArgConstStringChecker{value: "true"},
			&swagTagArgConstStringChecker{value: "false"},
		},
		errorMessage: "should be `true or false`",
	}
	swagTagArgHttpMethodUnionChecker = swagTagArgUnionChecker{
		options: []swagTagArgChecker{
			&swagTagArgConstStringChecker{value: "get", converter: strings.ToLower},
			&swagTagArgConstStringChecker{value: "post", converter: strings.ToLower},
			&swagTagArgConstStringChecker{value: "put", converter: strings.ToLower},
			&swagTagArgConstStringChecker{value: "patch", converter: strings.ToLower},
			&swagTagArgConstStringChecker{value: "delete", converter: strings.ToLower},
			&swagTagArgConstStringChecker{value: "head", converter: strings.ToLower},
			&swagTagArgConstStringChecker{value: "options", converter: strings.ToLower},
			&swagTagArgConstStringChecker{value: "trace", converter: strings.ToLower},
			&swagTagArgConstStringChecker{value: "connect", converter: strings.ToLower},
		},
		errorMessage: "should be `get, post, put, patch, delete, head, options, trace, or connect`",
	}
)

type swagTagArg interface {
	TagArgType() string
}

func NewSwagTagArg(def SwagTagArgDef, text string) (swagTagArg, error) {
	switch def.valueType {
	case swagTagArgDefTypeString:
		return &swagTagArgString{value: text}, nil
	case swagTagArgDefTypeGoType:
		return &swagTagArgGoType{value: text}, nil
	default:
		return nil, fmt.Errorf("unknown argDef.valueType: %d", def.valueType)
	}
}

type swagTagArgString struct {
	value string
}

func (s *swagTagArgString) TagArgType() string {
	return "string"
}

type swagTagArgGoType struct {
	value any
}

func (s *swagTagArgGoType) TagArgType() string {
	return "go-type"
}
