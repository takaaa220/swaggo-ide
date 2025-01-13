package swag

import (
	"strings"

	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/swag/parser"
	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/swag/tag"
)

type Attribute struct {
	Title       string
	Description string
}

func GetAttribute(line string) (*Attribute, error) {
	if !isCommentLine(line) {
		return nil, nil
	}

	firstToken, tokenizeArgs := parser.Tokenize(line)
	if !strings.HasPrefix(firstToken.Text, "@") {
		return nil, nil
	}
	if tokenizeArgs == nil {
		return nil, nil
	}
	swagTagDef := tag.NewSwagTagDef(strings.TrimPrefix(firstToken.Text, "@"))
	if !swagTagDef.IsValidTag() {
		return nil, nil
	}

	return &Attribute{
		Title:       swagTagDef.String(),
		Description: swagTagDef.Description,
	}, nil

}
