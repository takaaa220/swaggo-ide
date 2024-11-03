package swag

import (
	"fmt"
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

		_, checkErrors := check(line)
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

type checkError struct {
	message string
	start   int
	end     int
}

func check(line string) (bool, []checkError) {
	firstToken, tokenizeArgs := tokenize(line)
	if !strings.HasPrefix(firstToken.Text, "@") {
		return false, []checkError{}
	}
	if tokenizeArgs == nil {
		return false, []checkError{}
	}
	swagTagDef := newSwagTagDef(strings.TrimPrefix(firstToken.Text, "@"))
	if swagTagDef._type == swagTagTypeUnknown {
		return false, []checkError{}
	}

	checkErrors := []checkError{}
	i := 0
	for argToken := range tokenizeArgs(len(swagTagDef.args)) {
		def := swagTagDef.args[i]

		// TODO: move to tag.go
		text := trimBraces(argToken.Text)

		var arg swagTagArg
		switch def.valueType {
		case swagTagArgDefTypeString:
			arg = &swagTagArgString{value: text}
		case swagTagArgDefTypeGoType:
			arg = &swagTagArgGoType{value: text}
		default:
			panic(fmt.Errorf("unknown argDef.valueType: %d", def.valueType))
		}

		ok, errorMessages := def.check(arg)
		if !ok {
			checkErrors = append(checkErrors, checkError{
				message: strings.Join(errorMessages, ", "),
				start:   argToken.Start,
				end:     argToken.End,
			})
		}

		i++
	}

	if i < swagTagDef.requiredArgsCount {
		checkErrors = []checkError{{
			message: swagTagDef.errorMessage(),
			start:   firstToken.Start,
			end:     firstToken.End,
		}}
	}

	return len(checkErrors) == 0, checkErrors
}
