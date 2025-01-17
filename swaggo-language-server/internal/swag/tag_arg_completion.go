package swag

import (
	"strings"

	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/swag/parser"
	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/swag/tag"
)

type GetTagArgResult struct {
	Tag            *tag.SwagTagDef
	ActiveArgIndex int
}

func (r *GetTagArgResult) Candidates() []CompletionCandidate {
	candidates := r.Tag.Args[r.ActiveArgIndex].Candidates()

	res := make([]CompletionCandidate, len(candidates))
	for i, c := range candidates {
		res[i] = CompletionCandidate{
			NewText: c,
			Label:   c,
		}
	}

	return res
}

func GetTagArg(line string, position Position) (*GetTagArgResult, error) {
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

	i := -1
	positionIsLast := true
	for argToken := range tokenizeArgs(len(swagTagDef.Args)) {
		i++
		if int(position.Character) < argToken.End {
			if lastTokenEnd <= int(position.Character) && int(position.Character) < argToken.Start {
				return &GetTagArgResult{
					Tag:            &swagTagDef,
					ActiveArgIndex: i,
				}, nil
			}

			positionIsLast = false
			break
		}

		lastTokenEnd = argToken.End
	}
	if positionIsLast && i < len(swagTagDef.Args)-1 {
		return &GetTagArgResult{
			Tag:            &swagTagDef,
			ActiveArgIndex: i + 1,
		}, nil
	}

	return nil, nil
}
