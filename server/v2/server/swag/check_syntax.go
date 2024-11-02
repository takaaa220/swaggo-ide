package swag

import (
	"strings"

	"github.com/takaaa220/go-swag-ide/server/v2/server-sdk/protocol"
	"github.com/takaaa220/go-swag-ide/server/v2/server/util"
)

func CheckSyntax(uri string, src string) []protocol.Diagnostics {
	splitSrc := strings.Split(src, "\n")

	diagnostics := []protocol.Diagnostics{}
	for i, line := range splitSrc {
		if !util.IsCommentLine(line) {
			continue
		}

		trimmed, trimmedCount := util.TrimPrefixForComment(line)
		if !strings.HasPrefix(trimmed, "@") {
			continue
		}

		parser := NewSwagChecker()
		_, errors := parser.Check(trimmed)

		for _, errorMessage := range errors {
			diagnostics = append(diagnostics, protocol.Diagnostics{
				// TODO: Implement Range
				Range: protocol.Range{
					Start: protocol.Position{
						Line:      uint32(i),
						Character: uint32(trimmedCount),
					},
					End: protocol.Position{
						Line:      uint32(i),
						Character: uint32(trimmedCount + len(trimmed)),
					},
				},
				Severity: protocol.DiagnosticsSeverityError,
				Source:   "swag",
				Message:  errorMessage,
			})
		}
	}

	return diagnostics
}
