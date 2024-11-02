package swag

import (
	"strings"

	"github.com/takaaa220/go-swag-ide/server/v2/server-sdk/protocol"
)

func GetTagArgCompletionItems(line string, position protocol.Position) (*protocol.CompletionList, error) {
	swagTag, splitArgs := split(line)
	if !strings.HasPrefix(swagTag.Text, "@") || splitArgs == nil {
		return nil, nil
	}
	swagTagDef := newSwagTagDef(strings.TrimPrefix(swagTag.Text, "@"))
	if swagTagDef._type == swagTagTypeUnknown {
		return nil, nil
	}

	if swagTag.Start <= int(position.Character) && int(position.Character) < swagTag.End {
		return nil, nil
	}

	candidates := []string{}
	i := 0
	found := false
	for argSplitElement := range splitArgs(len(swagTagDef.args)) {
		if argSplitElement.Start <= int(position.Character) && int(position.Character) < argSplitElement.End {
			for _, argChecker := range swagTagDef.args[i].checkers {
				candidates = append(candidates, inner(argChecker)...)
			}

			found = true
			break
		}
		i++
	}

	if !found && i < len(swagTagDef.args)-1 {
		for _, argChecker := range swagTagDef.args[i].checkers {
			candidates = append(candidates, inner(argChecker)...)
		}
	}

	return nil, nil
}

func inner(checker swagTagArgChecker) []string {
	switch c := checker.(type) {
	case *swagTagArgUnionChecker:
		for _, option := range c.options {
			return inner(option)
		}
	case *swagTagArgConstStringChecker:
		return []string{c.value}
	}

	return []string{}
}
