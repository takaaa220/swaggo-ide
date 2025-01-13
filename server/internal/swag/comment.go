package swag

import (
	"strings"
)

func isCommentLine(line string) bool {
	return strings.HasPrefix(strings.TrimLeft(line, " \t"), "//")
}
