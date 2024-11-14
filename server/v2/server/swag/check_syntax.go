package swag

import (
	"strings"

	"github.com/takaaa220/go-swag-ide/server/v2/server-sdk/protocol"
	"github.com/takaaa220/go-swag-ide/server/v2/server/swag/parser"
	"github.com/takaaa220/go-swag-ide/server/v2/server/swag/tag"
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
	firstToken, tokenizeArgs := parser.Tokenize(line)
	if !strings.HasPrefix(firstToken.Text, "@") {
		return false, []checkError{}
	}
	if tokenizeArgs == nil {
		return false, []checkError{}
	}
	swagTagDef := tag.NewSwagTagDef(strings.TrimPrefix(firstToken.Text, "@"))
	if !swagTagDef.IsValidTag() {
		return false, []checkError{}
	}

	checkErrors := []checkError{}
	i := 0
	for argToken := range tokenizeArgs(len(swagTagDef.Args)) {
		def := swagTagDef.Args[i]

		arg, err := tag.NewSwagTagArg(def, argToken.Text)
		if err != nil {
			panic(err)
		}

		ok, errorMessages := def.Check(arg)
		if !ok {
			checkErrors = append(checkErrors, checkError{
				message: strings.Join(errorMessages, ", "),
				start:   argToken.Start,
				end:     argToken.End,
			})
		}

		i++
	}

	if i < swagTagDef.RequiredArgsCount() {
		checkErrors = []checkError{{
			message: swagTagDef.ErrorMessage(),
			start:   firstToken.Start,
			end:     firstToken.End,
		}}
	}

	return len(checkErrors) == 0, checkErrors
}
