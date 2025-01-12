package swag

import (
	"strings"

	"github.com/takaaa220/swaggo-ide/server/internal/handler/protocol"
	"github.com/takaaa220/swaggo-ide/server/internal/handler/swag/parser"
	"github.com/takaaa220/swaggo-ide/server/internal/handler/swag/tag"
	"github.com/takaaa220/swaggo-ide/server/internal/handler/util"
)

type Attribute struct {
	Title       string
	Description string
}

func GetAttribute(line string, position protocol.Position) (*Attribute, error) {
	if !util.IsCommentLine(line) {
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
