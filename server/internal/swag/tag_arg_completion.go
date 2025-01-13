package swag

import (
	"strings"

	"github.com/takaaa220/swaggo-ide/server/internal/swag/parser"
	"github.com/takaaa220/swaggo-ide/server/internal/swag/tag"
)

func GetTagArgCompletionItems(line string, position Position) ([]CompletionCandidate, error) {
	if !isCommentLine(line) {
		return nil, nil
	}

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

	completionCandidates := make([]CompletionCandidate, len(candidates))
	for i, candidate := range candidates {
		completionCandidates[i] = CompletionCandidate{
			Label:   candidate,
			NewText: candidate,
		}
	}

	return completionCandidates, nil
}
