package swag

import (
	"sort"

	"github.com/takaaa220/go-swag-ide/server/v2/server-sdk/protocol"
)

func GetCompletionItems(position protocol.Position) (*protocol.CompletionList, error) {
	completionItems := make([]protocol.CompletionItem, len(swagTags))
	for i, tag := range swagTags {
		completionItems[i] = protocol.CompletionItem{
			Label: tag._type.String(),
			Kind:  protocol.CompletionItemKindKeyword,
			TextEdit: &protocol.TextEdit{
				NewText: tag.String(),
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

	sort.Slice(completionItems, func(i, j int) bool {
		return completionItems[i].Label < completionItems[j].Label
	})

	return &protocol.CompletionList{
		IsIncomplete: false,
		Items:        completionItems,
	}, nil
}
