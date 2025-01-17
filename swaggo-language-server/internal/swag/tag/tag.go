package tag

import (
	"fmt"
	"strings"
)

type SwagTagDef struct {
	Type              swagTagType
	Args              []SwagTagArgDef
	Description       string
	requiredArgsCount int
}

func NewSwagTagDef(tagTypeString string) SwagTagDef {
	for _, tag := range SwagTags {
		if strings.ToLower(tagTypeString) == strings.ToLower(string(tag.Type)) {
			return tag
		}
	}

	return swagTagUnknown
}

func (s *SwagTagDef) IsValidTag() bool {
	return s.Type != swagTagTypeUnknown
}

func (s *SwagTagDef) RequiredArgsCount() int {
	return s.requiredArgsCount
}

func (s *SwagTagDef) ErrorMessage() string {
	return fmt.Sprintf("Should be `%s`.", s.String())
}

func (s *SwagTagDef) String() string {
	argsNameList := make([]string, s.requiredArgsCount)
	for i := range s.requiredArgsCount {
		argsNameList[i] = s.Args[i].label
	}

	return fmt.Sprintf("%s %s", s.Type, strings.Join(argsNameList, " "))
}

var (
	swagTagSummary = SwagTagDef{
		Type:              swagTagTypeSummary,
		Args:              []SwagTagArgDef{newSwagTagStringArgDef("SUMMARY")},
		Description:       "A short summary of the operation.",
		requiredArgsCount: 1,
	}
	swagTagDescription = SwagTagDef{
		Type:              swagTagTypeDescription,
		Args:              []SwagTagArgDef{newSwagTagStringArgDef("DESCRIPTION")},
		Description:       "A verbose explanation of the operation.",
		requiredArgsCount: 1,
	}
	swagTagTags = SwagTagDef{
		Type:              swagTagTypeTags,
		Args:              []SwagTagArgDef{newSwagTagStringArgDef("TAG1,TAG2")},
		Description:       "A list of tags for API documentation control.",
		requiredArgsCount: 1,
	}
	swagTagAccept = SwagTagDef{
		Type:              swagTagTypeAccept,
		Args:              []SwagTagArgDef{newSwagTagUnionArgDef("MIME_TYPE", swagTagArgMimeTypeUnionChecker)},
		Description:       "A list of MIME types the operation can consume.",
		requiredArgsCount: 1,
	}
	swagTagProduce = SwagTagDef{
		Type:              swagTagTypeProduce,
		Args:              []SwagTagArgDef{newSwagTagUnionArgDef("MIME_TYPE", swagTagArgMimeTypeUnionChecker)},
		Description:       "A list of MIME types the operation can produce.",
		requiredArgsCount: 1,
	}
	swagTagParam = SwagTagDef{
		Type: swagTagTypeParam,
		Args: []SwagTagArgDef{
			newSwagTagStringArgDef("PARAM_NAME"),
			newSwagTagUnionArgDef("PARAM_TYPE", swagTagArgParamTypeUnionChecker),
			newSwagTagUnionArgDef("GO_TYPE", swagTagArgGoDataTypeUnionChecker),
			newSwagTagBoolArgDef("REQUIRED"),
			wrapArgDefWithBraces('"', newSwagTagStringArgDef("DESCRIPTION")),
			newSwagTagStringArgDef("ATTRIBUTE"),
		},
		Description:       "Describes a single operation parameter.",
		requiredArgsCount: 5,
	}
	swagTagSuccess = SwagTagDef{
		Type: swagTagTypeSuccess,
		Args: []SwagTagArgDef{
			newSwagTagIntArgDef("STATUS_CODE"),
			wrapArgDefWithBraces('{', newSwagTagUnionArgDef("DATA_TYPE", swagTagArgDataTypeUnionChecker)),
			newSwagTagUnionArgDef("GO_TYPE", swagTagArgGoDataTypeUnionChecker),
			newSwagTagStringArgDef("DESCRIPTION"),
		},
		Description:       "A success response.",
		requiredArgsCount: 3,
	}
	swagTagFailure = SwagTagDef{
		Type: swagTagTypeFailure,
		Args: []SwagTagArgDef{
			newSwagTagIntArgDef("STATUS_CODE"),
			wrapArgDefWithBraces('{', newSwagTagUnionArgDef("DATA_TYPE", swagTagArgDataTypeUnionChecker)),
			newSwagTagUnionArgDef("GO_TYPE", swagTagArgGoDataTypeUnionChecker),
			newSwagTagStringArgDef("DESCRIPTION"),
		},
		Description:       "A failure response.",
		requiredArgsCount: 3,
	}
	swagTagRouter = SwagTagDef{
		Type: swagTagTypeRouter,
		Args: []SwagTagArgDef{
			newSwagTagStringArgDef("PATH"),
			wrapArgDefWithBraces('[', newSwagTagUnionArgDef("HTTP_METHOD", &swagTagArgHttpMethodUnionChecker)),
		},
		Description:       "A router definition. This is required for the operation.",
		requiredArgsCount: 2,
	}
	swagTagID = SwagTagDef{
		Type: swagTagTypeID,
		Args: []SwagTagArgDef{
			newSwagTagStringArgDef("ID"),
		},
		Description:       "A unique identifier for the operation.",
		requiredArgsCount: 1,
	}
	swagTagHeader = SwagTagDef{
		Type: swagTagTypeHeader,
		Args: []SwagTagArgDef{
			newSwagTagIntArgDef("STATUS_CODE"),
			wrapArgDefWithBraces('{', newSwagTagUnionArgDef("DATA_TYPE", swagTagArgDataTypeUnionChecker)),
			newSwagTagStringArgDef("HEADER_NAME"),
			newSwagTagStringArgDef("COMMENT"),
		},
		Description:       "A header definition.",
		requiredArgsCount: 4,
	}
	swagTagUnknown = SwagTagDef{
		Type: swagTagTypeUnknown,
	}
	SwagTags = []SwagTagDef{
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
