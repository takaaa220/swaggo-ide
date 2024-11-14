package swag

import (
	"strings"

	"github.com/takaaa220/go-swag-ide/server/internal/server-sdk/protocol"
	"github.com/takaaa220/go-swag-ide/server/internal/server/swag/parser"
	"github.com/takaaa220/go-swag-ide/server/internal/server/swag/tag"
)

func GetTagArgCompletionItems(line string, position protocol.Position) (*protocol.CompletionList, error) {
	firstToken, tokenizeArgs := parser.Tokenize(line)
	if !strings.HasPrefix(firstToken.Text, "@") || tokenizeArgs == nil {
		return nil, nil
	}
	swagTagDef := tag.NewSwagTagDef(strings.TrimPrefix(firstToken.Text, "@"))
	if !swagTagDef.IsValidTag() {
		return nil, nil
	}

	if int(position.Character) < firstToken.End {
		return nil, nil
	}

	lastTokenEnd := firstToken.End

	candidates := []string{}
	i := -1
	positionIsLast := true
	for argToken := range tokenizeArgs(len(swagTagDef.Args)) {
		i++
		if int(position.Character) < argToken.End {
			if lastTokenEnd <= int(position.Character) && int(position.Character) < argToken.Start {
				candidates = append(candidates, swagTagDef.Args[i].Candidates()...)
			}

			positionIsLast = false
			break
		}

		lastTokenEnd = argToken.End
	}

	if positionIsLast && i < len(swagTagDef.Args)-1 {
		candidates = append(candidates, swagTagDef.Args[i+1].Candidates()...)
	}

	completionItems := make([]protocol.CompletionItem, len(candidates))
	for i, candidate := range candidates {
		completionItems[i] = protocol.CompletionItem{
			Label: candidate,
			Kind:  protocol.CompletionItemKindKeyword,
			TextEdit: protocol.TextEdit{
				Range: protocol.Range{
					Start: protocol.Position{Line: position.Line, Character: position.Character},
					End:   protocol.Position{Line: position.Line, Character: position.Character},
				},
				NewText: candidate,
			},
		}
	}

	return &protocol.CompletionList{
		IsIncomplete: false,
		Items:        completionItems,
	}, nil
}
