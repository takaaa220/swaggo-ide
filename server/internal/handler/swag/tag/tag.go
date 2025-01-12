package tag

import (
	"fmt"
	"strings"
)

type swagTagDef struct {
	Type              swagTagType
	Args              []swagTagArgDef
	requiredArgsCount int
}

func NewSwagTagDef(tagTypeString string) swagTagDef {
	for _, tag := range SwagTags {
		if strings.ToLower(tagTypeString) == strings.ToLower(string(tag.Type)) {
			return tag
		}
	}

	return swagTagUnknown
}

func (s *swagTagDef) IsValidTag() bool {
	return s.Type != swagTagTypeUnknown
}

func (s *swagTagDef) RequiredArgsCount() int {
	return s.requiredArgsCount
}

func (s *swagTagDef) ErrorMessage() string {
	return fmt.Sprintf("Should be `%s`.", s.String())
}

func (s *swagTagDef) String() string {
	argsNameList := make([]string, s.requiredArgsCount)
	for i := range s.requiredArgsCount {
		argsNameList[i] = s.Args[i].label
	}

	return fmt.Sprintf("%s %s", s.Type, strings.Join(argsNameList, " "))
}

var (
	swagTagSummary = swagTagDef{
		Type:              swagTagTypeSummary,
		Args:              []swagTagArgDef{newSwagTagStringArgDef("SUMMARY")},
		requiredArgsCount: 1,
	}
	swagTagDescription = swagTagDef{
		Type:              swagTagTypeDescription,
		Args:              []swagTagArgDef{newSwagTagStringArgDef("DESCRIPTION")},
		requiredArgsCount: 1,
	}
	swagTagTags = swagTagDef{
		Type:              swagTagTypeTags,
		Args:              []swagTagArgDef{newSwagTagStringArgDef("TAG1,TAG2")},
		requiredArgsCount: 1,
	}
	swagTagAccept = swagTagDef{
		Type:              swagTagTypeAccept,
		Args:              []swagTagArgDef{newSwagTagUnionArgDef("MIME_TYPE", swagTagArgMimeTypeUnionChecker)},
		requiredArgsCount: 1,
	}
	swagTagProduce = swagTagDef{
		Type:              swagTagTypeProduce,
		Args:              []swagTagArgDef{newSwagTagUnionArgDef("MIME_TYPE", swagTagArgMimeTypeUnionChecker)},
		requiredArgsCount: 1,
	}
	swagTagParam = swagTagDef{
		Type: swagTagTypeParam,
		Args: []swagTagArgDef{
			newSwagTagStringArgDef("PARAM_NAME"),
			newSwagTagUnionArgDef("PARAM_TYPE", swagTagArgParamTypeUnionChecker),
			newSwagTagUnionArgDef("GO_TYPE", swagTagArgGoDataTypeUnionChecker),
			newSwagTagBoolArgDef("REQUIRED"),
			wrapArgDefWithBraces('"', newSwagTagStringArgDef("DESCRIPTION")),
			newSwagTagStringArgDef("ATTRIBUTE"),
		},
		requiredArgsCount: 5,
	}
	swagTagSuccess = swagTagDef{
		Type: swagTagTypeSuccess,
		Args: []swagTagArgDef{
			newSwagTagIntArgDef("STATUS_CODE"),
			wrapArgDefWithBraces('{', newSwagTagUnionArgDef("DATA_TYPE", swagTagArgDataTypeUnionChecker)),
			newSwagTagUnionArgDef("GO_TYPE", swagTagArgGoDataTypeUnionChecker),
			newSwagTagStringArgDef("DESCRIPTION"),
		},
		requiredArgsCount: 3,
	}
	swagTagFailure = swagTagDef{
		Type: swagTagTypeFailure,
		Args: []swagTagArgDef{
			newSwagTagIntArgDef("STATUS_CODE"),
			wrapArgDefWithBraces('{', newSwagTagUnionArgDef("DATA_TYPE", swagTagArgDataTypeUnionChecker)),
			newSwagTagUnionArgDef("GO_TYPE", swagTagArgGoDataTypeUnionChecker),
			newSwagTagStringArgDef("DESCRIPTION"),
		},
		requiredArgsCount: 3,
	}
	swagTagRouter = swagTagDef{
		Type: swagTagTypeRouter,
		Args: []swagTagArgDef{
			newSwagTagStringArgDef("PATH"),
			wrapArgDefWithBraces('[', newSwagTagUnionArgDef("HTTP_METHOD", &swagTagArgHttpMethodUnionChecker)),
		},
		requiredArgsCount: 2,
	}
	swagTagID = swagTagDef{
		Type: swagTagTypeID,
		Args: []swagTagArgDef{
			newSwagTagStringArgDef("ID"),
		},
		requiredArgsCount: 1,
	}
	swagTagHeader = swagTagDef{
		Type: swagTagTypeHeader,
		Args: []swagTagArgDef{
			newSwagTagIntArgDef("STATUS_CODE"),
			wrapArgDefWithBraces('{', newSwagTagUnionArgDef("DATA_TYPE", swagTagArgDataTypeUnionChecker)),
			newSwagTagStringArgDef("HEADER_NAME"),
			newSwagTagStringArgDef("COMMENT"),
		},
		requiredArgsCount: 4,
	}
	swagTagUnknown = swagTagDef{
		Type: swagTagTypeUnknown,
	}
	SwagTags = []swagTagDef{
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

func (s swagTagType) String() string {
	return "@" + string(s)
}

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

func (s swagTagType) IsRouter() bool {
	return s == swagTagTypeRouter
}
