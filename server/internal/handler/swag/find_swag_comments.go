package swag

import (
	"strings"

	"github.com/takaaa220/swaggo-ide/server/internal/handler/swag/parser"
	"github.com/takaaa220/swaggo-ide/server/internal/handler/swag/tag"
	"github.com/takaaa220/swaggo-ide/server/internal/handler/util"
)

type SwagCommentsRange struct {
	// Start is the start line index of the swag comments.
	Start int
	// End is the end line index of the swag comments.
	End int
}

func FindSwagComments(src string) []SwagCommentsRange {
	splitSrc := strings.Split(src, "\n")

	var (
		swagCommentsLineIndexes []SwagCommentsRange
		commentsStartIndex      *int
		isSwagComment           bool
	)
	for i, line := range splitSrc {
		if !util.IsCommentLine(line) {
			if commentsStartIndex != nil {
				if isSwagComment {
					swagCommentsLineIndexes = append(swagCommentsLineIndexes, SwagCommentsRange{
						Start: *commentsStartIndex,
						End:   i - 1,
					})
				}
				commentsStartIndex = nil
				isSwagComment = false
			}

			continue
		}

		if commentsStartIndex == nil {
			commentsStartIndex = &i
		}

		firstToken, tokenizeArgs := parser.Tokenize(line)
		if !strings.HasPrefix(firstToken.Text, "@") {
			continue
		}
		if tokenizeArgs == nil {
			continue
		}
		swagTagDef := tag.NewSwagTagDef(strings.TrimPrefix(firstToken.Text, "@"))
		if !swagTagDef.IsValidTag() {
			continue
		}

		isSwagComment = true
	}

	if commentsStartIndex != nil && isSwagComment {
		swagCommentsLineIndexes = append(swagCommentsLineIndexes, SwagCommentsRange{
			Start: *commentsStartIndex,
			End:   len(splitSrc) - 1,
		})
	}

	return swagCommentsLineIndexes
}
