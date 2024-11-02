package server

import (
	"strings"

	"github.com/takaaa220/go-swag-ide/server/v2/server-sdk/protocol"
)

func isInFunctionComment(src string, pos protocol.Position) (bool, error) {
	splitSrc := strings.Split(src, "\n")
	if int(pos.Line) >= len(splitSrc) {
		return false, nil
	}

	line := splitSrc[pos.Line]

	return strings.HasPrefix(strings.TrimLeft(line, " \t"), "//"), nil
}
