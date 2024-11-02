package util

import (
	"strings"

	"github.com/takaaa220/go-swag-ide/server/v2/server-sdk/protocol"
)

func IsInComment(src string, pos protocol.Position) bool {
	// TODO: `[POS_IS_HERE] // comment` ← should return false

	splitSrc := strings.Split(src, "\n")
	if len(splitSrc) <= int(pos.Line) {
		return false
	}

	return IsCommentLine(splitSrc[pos.Line])
}

func IsCommentLine(line string) bool {
	return strings.HasPrefix(strings.TrimLeft(line, " \t"), "//")
}

// TrimPrefixForComment trims the prefix of the line if it is a comment line.
//
//	input: `  // comment comment`
//	output: `comment comment`
func TrimPrefixForComment(line string) (string, int) {
	if !IsCommentLine(line) {
		return line, 0
	}

	trimmedCount := 0
	for _, c := range line {
		if c == ' ' || c == '\t' || c == '/' {
			trimmedCount++
			continue
		}

		break
	}

	return line[trimmedCount:], trimmedCount
}
