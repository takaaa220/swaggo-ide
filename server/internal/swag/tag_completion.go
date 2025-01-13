package swag

import (
	"sort"

	"github.com/takaaa220/swaggo-ide/server/internal/swag/tag"
)

func GetTagCompletionItems(line string) ([]CompletionCandidate, error) {
	if !isCommentLine(line) {
		return nil, nil
	}

	candidates := make([]CompletionCandidate, len(tag.SwagTags))
	for i, tag := range tag.SwagTags {
		candidates[i] = CompletionCandidate{
			Label:   tag.Type.String(),
			NewText: tag.String(),
		}
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Label < candidates[j].Label
	})

	return candidates, nil
}
