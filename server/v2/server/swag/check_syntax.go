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

		parser := NewSwagChecker()
		_, checkErrors := parser.Check(line)

		for _, checkError := range checkErrors {
			diagnostics = append(diagnostics, protocol.Diagnostics{
				Range: protocol.Range{
					Start: protocol.Position{
						Line:      uint32(i),
						Character: uint32(checkError.start),
					},
					End: protocol.Position{
						Line:      uint32(i),
						Character: uint32(checkError.end),
					},
				},
				Severity: protocol.DiagnosticsSeverityError,
				Source:   "swag",
				Message:  checkError.message,
			})
		}
	}

	return diagnostics
}
