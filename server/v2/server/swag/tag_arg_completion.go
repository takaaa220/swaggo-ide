package swag

import (
	"strings"

	"github.com/takaaa220/go-swag-ide/server/v2/server-sdk/protocol"
)

func GetTagArgCompletionItems(line string, position protocol.Position) (*protocol.CompletionList, error) {
	firstToken, tokenizeArgs := tokenize(line)
	if !strings.HasPrefix(firstToken.Text, "@") || tokenizeArgs == nil {
		return nil, nil
	}
	swagTagDef := newSwagTagDef(strings.TrimPrefix(firstToken.Text, "@"))
	if swagTagDef._type == swagTagTypeUnknown {
		return nil, nil
	}

	if int(position.Character) < firstToken.End {
		return nil, nil
	}

	lastTokenEnd := firstToken.End

	candidates := []string{}
	i := -1
	positionIsLast := true
	for argToken := range tokenizeArgs(len(swagTagDef.args)) {
		i++
		if int(position.Character) < argToken.End {
			if lastTokenEnd <= int(position.Character) && int(position.Character) < argToken.Start {
				for _, argChecker := range swagTagDef.args[i].checkers {
					candidates = append(candidates, inner(argChecker)...)
				}
			}

			positionIsLast = false
			break
		}

		lastTokenEnd = argToken.End
	}

	if positionIsLast && i < len(swagTagDef.args)-1 {
		for _, argChecker := range swagTagDef.args[i+1].checkers {
			candidates = append(candidates, inner(argChecker)...)
		}
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

func inner(checker swagTagArgChecker) []string {
	candidates := []string{}
	switch c := checker.(type) {
	case *swagTagArgUnionChecker:
		for _, option := range c.options {
			candidates = append(candidates, inner(option)...)
		}
	case *swagTagArgConstStringChecker:
		candidates = append(candidates, c.value)
	}

	return candidates
}
