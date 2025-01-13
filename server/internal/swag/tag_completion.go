package swag

import (
	"sort"

	"github.com/takaaa220/swaggo-ide/server/internal/handler/protocol"
	"github.com/takaaa220/swaggo-ide/server/internal/swag/tag"
)

func GetTagCompletionItems(line string, position protocol.Position) (*protocol.CompletionList, error) {
	if !isCommentLine(line) {
		return nil, nil
	}

	completionItems := make([]protocol.CompletionItem, len(tag.SwagTags))
	for i, tag := range tag.SwagTags {
		completionItems[i] = protocol.CompletionItem{
			Label: tag.Type.String(),
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
