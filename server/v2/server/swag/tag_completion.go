package swag

import "github.com/takaaa220/go-swag-ide/server/v2/server-sdk/protocol"

type swagTagOld struct {
	label string
	args  []string
}

var tags = []swagTagOld{
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

func GetCompletionItems(position protocol.Position) (*protocol.CompletionList, error) {
	kind := protocol.CompletionItemKindKeyword
	completionItems := make([]protocol.CompletionItem, len(tags))
	for i, tag := range tags {
		completionText := tag.label
		for _, arg := range tag.args {
			completionText += "  " + arg
		}

		completionItems[i] = protocol.CompletionItem{
			Label: tag.label,
			Kind:  kind,
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
