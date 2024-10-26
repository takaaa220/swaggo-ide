package internal

import (
	"go/token"
	"net/url"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

type completionCreator interface {
	getCompletionItems(position protocol.Position) (*protocol.CompletionList, error)
}

func completionCreatorFactory(params *protocol.CompletionParams) completionCreator {
	parsedUri, err := url.Parse(params.TextDocument.URI)
	if err != nil {
		return &noopCompletionCreator{}
	}

	t, err := isInFunctionComment(token.Position{
		Filename: parsedUri.Path,
		Line:     int(params.Position.Line + 1),
		Column:   int(params.Position.Character + 1),
		Offset:   -1,
	})
	if err != nil {
		return &noopCompletionCreator{}
	}
	if !t {
		return &noopCompletionCreator{}
	}

	if *params.Context.TriggerCharacter == "@" {
		return &tagCompletionCreator{}
	}

	return &noopCompletionCreator{}
}

type noopCompletionCreator struct{}

func (ncc *noopCompletionCreator) getCompletionItems(position protocol.Position) (*protocol.CompletionList, error) {
	return nil, nil
}

type tagCompletionCreator struct{}

type swagTag struct {
	label string
	args  []string
}

var swagTags = []swagTag{
	{
		label: "@Summary",
		args:  []string{"SUMMARY"},
	},
	{
		label: "@Description",
		args:  []string{"DESCRIPTION"},
	},
	{
		label: "@Tags",
		args:  []string{"TAG1,TAG2"},
	},
	{
		label: "@Accept",
		args:  []string{"MIME_TYPE"},
	},
	{
		label: "@Produce",
		args:  []string{"MIME_TYPE"},
	},
	{
		label: "@Param",
		args:  []string{"PARAM_NAME", "PARAM_TYPE", "DATA_TYPE", "REQUIRED(bool)", "DESCRIPTION", "ATTRIBUTE(optional)"},
	},
	{
		label: "@Success",
		args:  []string{"STATUS_CODE", "{DATA_TYPE}", "DESCRIPTION"},
	},
	{
		label: "@Failure",
		args:  []string{"STATUS_CODE", "{DATA_TYPE}", "DESCRIPTION"},
	},
	{
		label: "@Router",
		args:  []string{"PATH", "[METHOD]"},
	},
	{
		label: "@Security",
		args:  []string{},
	},
	{
		label: "@ID",
		args:  []string{"ID"},
	},
	{
		label: "@Header",
		args:  []string{"STATUS_CODE", "{PARAM_TYPE}", "DATA_TYPE", "COMMENT"},
	},
}

func (tcc *tagCompletionCreator) getCompletionItems(position protocol.Position) (*protocol.CompletionList, error) {
	kind := protocol.CompletionItemKindKeyword
	completionItems := make([]protocol.CompletionItem, len(swagTags))
	for i, tag := range swagTags {
		completionText := tag.label
		for _, arg := range tag.args {
			completionText += "  " + arg
		}

		completionItems[i] = protocol.CompletionItem{
			Label: tag.label,
			Kind:  &kind,
			TextEdit: &protocol.TextEdit{
				NewText: completionText,
				Range: protocol.Range{
					Start: protocol.Position{
						Line:      position.Line,
						Character: position.Character - 1,
					},
					End: position,
				},
			},
		}

	}

	return &protocol.CompletionList{
		IsIncomplete: false,
		Items:        completionItems,
	}, nil
}
