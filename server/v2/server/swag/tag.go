package swag

import (
	"fmt"
	"strconv"
	"strings"
)

type swagTagDef struct {
	_type             swagTagType
	args              []swagTagArgDef
	requiredArgsCount int
}

func newSwagTagDef(tagTypeString string) swagTagDef {
	for _, tag := range swagTags {
		if strings.ToLower(tagTypeString) == strings.ToLower(string(tag._type)) {
			return tag
		}
	}

	return swagTagUnknown
}

func (s *swagTagDef) errorMessage() string {
	argsNameList := make([]string, s.requiredArgsCount)
	for i := range s.requiredArgsCount {
		argsNameList[i] = s.args[i].label
	}

	return fmt.Sprintf("Should be `@%s %s`.", s._type, strings.Join(argsNameList, ", "))
}

var (
	swagTagSummary = swagTagDef{
		_type:             swagTagTypeSummary,
		args:              []swagTagArgDef{newSwagTagStringArgDef("SUMMARY")},
		requiredArgsCount: 1,
	}
	swagTagDescription = swagTagDef{
		_type:             swagTagTypeDescription,
		args:              []swagTagArgDef{newSwagTagStringArgDef("DESCRIPTION")},
		requiredArgsCount: 1,
	}
	swagTagTags = swagTagDef{
		_type:             swagTagTypeTags,
		args:              []swagTagArgDef{newSwagTagStringArgDef("TAG1,TAG2")},
		requiredArgsCount: 1,
	}
	swagTagAccept = swagTagDef{
		_type:             swagTagTypeAccept,
		args:              []swagTagArgDef{newSwagTagUnionArgDef("MIME_TYPE", &swagTagArgMimeTypeUnionChecker)},
		requiredArgsCount: 1,
	}
	swagTagProduce = swagTagDef{
		_type:             swagTagTypeProduce,
		args:              []swagTagArgDef{newSwagTagUnionArgDef("MIME_TYPE", &swagTagArgMimeTypeUnionChecker)},
		requiredArgsCount: 1,
	}
	swagTagParam = swagTagDef{
		_type: swagTagTypeParam,
		args: []swagTagArgDef{
			newSwagTagStringArgDef("PARAM_NAME"),
			newSwagTagUnionArgDef("PARAM_TYPE", &swagTagArgParamTypeUnionChecker),
			newSwagTagUnionArgDef("DATA_TYPE", &swagTagArgDataTypeUnionChecker),
			newSwagTagBoolArgDef("REQUIRED"),
			newSwagTagStringArgDef("DESCRIPTION"),
			newSwagTagStringArgDef("ATTRIBUTE"),
		},
		requiredArgsCount: 5,
	}
	swagTagSuccess = swagTagDef{
		_type: swagTagTypeSuccess,
		args: []swagTagArgDef{
			newSwagTagIntArgDef("STATUS_CODE"),
			newSwagTagUnionArgDef("{PARAM_TYPE}", &swagTagArgParamTypeUnionChecker),
			newSwagTagUnionArgDef("DATA_TYPE", &swagTagArgDataTypeUnionChecker),
			newSwagTagStringArgDef("DESCRIPTION"),
		},
		requiredArgsCount: 3,
	}
	swagTagFailure = swagTagDef{
		_type: swagTagTypeFailure,
		args: []swagTagArgDef{
			newSwagTagIntArgDef("STATUS_CODE"),
			newSwagTagUnionArgDef("{PARAM_TYPE}", &swagTagArgParamTypeUnionChecker),
			newSwagTagUnionArgDef("DATA_TYPE", &swagTagArgDataTypeUnionChecker),
			newSwagTagStringArgDef("DESCRIPTION"),
		},
		requiredArgsCount: 3,
	}
	swagTagRouter = swagTagDef{
		_type: swagTagTypeRouter,
		args: []swagTagArgDef{
			newSwagTagStringArgDef("PATH"),
			newSwagTagUnionArgDef("[HTTP_METHOD]", &swagTagArgHttpMethodUnionChecker),
		},
	}
	swagTagID = swagTagDef{
		_type: swagTagTypeID,
		args: []swagTagArgDef{
			newSwagTagStringArgDef("ID"),
		},
		requiredArgsCount: 1,
	}
	swagTagHeader = swagTagDef{
		_type: swagTagTypeHeader,
		args: []swagTagArgDef{
			newSwagTagIntArgDef("STATUS_CODE"),
			// TODO: Check
			// newSwagTagUnionArgDef("PARAM_TYPE", &swagTagArgParamTypeUnionChecker),
			newSwagTagStringArgDef("{PARAM_TYPE}"),
			newSwagTagUnionArgDef("HEADER_NAME", &swagTagArgDataTypeUnionChecker),
			newSwagTagStringArgDef("COMMENT"),
		},
		requiredArgsCount: 4,
	}
	swagTagUnknown = swagTagDef{
		_type: swagTagTypeUnknown,
	}
	swagTags = []swagTagDef{
		swagTagSummary,
		swagTagDescription,
		swagTagTags,
		swagTagAccept,
		swagTagProduce,
		swagTagParam,
		swagTagSuccess,
		swagTagFailure,
		swagTagRouter,
		swagTagID,
		swagTagHeader,
	}
)

type swagTagType string

const (
	swagTagTypeSummary     swagTagType = "Summary"
	swagTagTypeDescription swagTagType = "Description"
	swagTagTypeTags        swagTagType = "Tags"
	swagTagTypeAccept      swagTagType = "Accept"
	swagTagTypeProduce     swagTagType = "Produce"
	swagTagTypeParam       swagTagType = "Param"
	swagTagTypeSuccess     swagTagType = "Success"
	swagTagTypeFailure     swagTagType = "Failure"
	swagTagTypeRouter      swagTagType = "Router"
	swagTagTypeID          swagTagType = "ID"
	swagTagTypeHeader      swagTagType = "Header"
	swagTagTypeUnknown     swagTagType = "-"
)

type swagTagArgDefType int

const (
	swagTagArgDefTypeString swagTagArgDefType = iota
	swagTagArgDefTypeGoType
)

type swagTagArgDef struct {
	valueType swagTagArgDefType
	label     string
	checkers  []swagTagArgChecker
}

func newSwagTagStringArgDef(label string, optionalCheckers ...swagTagArgChecker) swagTagArgDef {
	return swagTagArgDef{
		valueType: swagTagArgDefTypeString,
		label:     label,
		checkers:  append([]swagTagArgChecker{&swagTagArgStringChecker{}}, optionalCheckers...),
	}
}

func newSwagTagIntArgDef(label string, optionalCheckers ...swagTagArgChecker) swagTagArgDef {
	return swagTagArgDef{
		valueType: swagTagArgDefTypeString,
		label:     label,
		checkers:  append([]swagTagArgChecker{&swagTagArgIntChecker{}}, optionalCheckers...),
	}
}

func newSwagTagBoolArgDef(label string, optionalCheckers ...swagTagArgChecker) swagTagArgDef {
	return swagTagArgDef{
		valueType: swagTagArgDefTypeString,
		label:     label,
		checkers:  append([]swagTagArgChecker{&swagTagArgBoolChecker}, optionalCheckers...),
	}
}

func newSwagTagGoTypeArgDef(label string, optionalCheckers ...swagTagArgChecker) swagTagArgDef {
	return swagTagArgDef{
		valueType: swagTagArgDefTypeGoType,
		label:     label,
		checkers:  append([]swagTagArgChecker{&swagTagArgGoTypeChecker{}}, optionalCheckers...),
	}
}

func newSwagTagConstStringArgDef(label string, value string, optionalCheckers ...swagTagArgChecker) swagTagArgDef {
	return swagTagArgDef{
		valueType: swagTagArgDefTypeString,
		label:     label,
		checkers:  append([]swagTagArgChecker{&swagTagArgConstStringChecker{value: value}}, optionalCheckers...),
	}
}

func newSwagTagUnionArgDef(label string, checker *swagTagArgUnionChecker, optionalCheckers ...swagTagArgChecker) swagTagArgDef {
	return swagTagArgDef{
		valueType: swagTagArgDefTypeString,
		label:     label,
		checkers:  append([]swagTagArgChecker{checker}, optionalCheckers...),
	}
}

func (s *swagTagArgDef) isValid(arg swagTagArg) (bool, []string) {
	errors := []string{}
	for _, checker := range s.checkers {
		if ok, errorMessage := checker.check(arg, *s); !ok {
			errors = append(errors, fmt.Sprintf("%s %s.", s.label, errorMessage))
		}
	}

	return len(errors) == 0, errors
}

type swagTagArgChecker interface {
	// check returns check result and error message
	check(arg swagTagArg, def swagTagArgDef) (bool, string)
}

type swagTagArgStringChecker struct{}

func (s *swagTagArgStringChecker) check(arg swagTagArg, def swagTagArgDef) (bool, string) {
	if _, ok := arg.(*swagTagArgString); !ok {
		return false, "invalid string"
	}

	return true, ""
}

type swagTagArgIntChecker struct{}

func (s *swagTagArgIntChecker) check(arg swagTagArg, def swagTagArgDef) (bool, string) {
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

type swagTagArgConstStringChecker struct {
	value     string
	converter func(string) string
}

func (s *swagTagArgConstStringChecker) check(arg swagTagArg, def swagTagArgDef) (bool, string) {
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

type swagTagArgUnionChecker struct {
	options      []swagTagArgChecker
	errorMessage string
}

func (s *swagTagArgUnionChecker) check(arg swagTagArg, def swagTagArgDef) (bool, string) {
	for _, option := range s.options {
		if ok, _ := option.check(arg, def); ok {
			return true, ""
		}
	}

	return false, s.errorMessage
}

type swagTagArgGoTypeChecker struct{}

func (s *swagTagArgGoTypeChecker) check(arg swagTagArg, def swagTagArgDef) (bool, string) {
	// TODO: Implement
	return true, ""
}

var (
	userDefinedType                = &swagTagArgGoTypeChecker{}
	swagTagArgMimeTypeUnionChecker = swagTagArgUnionChecker{
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
	swagTagArgParamTypeUnionChecker = swagTagArgUnionChecker{
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
	swagTagArgDataTypeUnionChecker = swagTagArgUnionChecker{
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

type swagTagArgInt struct {
	value int
}

func (s *swagTagArgInt) TagArgType() string {
	return "int"
}
