package swag

import (
	"strings"

	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/swag/parser"
	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/swag/tag"
)

type SyntaxError struct {
	Range   Range
	Message string
}

func CheckSyntax(uri string, src string) []SyntaxError {
	splitSrc := strings.Split(src, "\n")

	var checkers []*apiRouteInfoChecker
	var checker *apiRouteInfoChecker
	for i, line := range splitSrc {
		if !isCommentLine(line) {
			if checker != nil {
				checkers = append(checkers, checker)
				checker = nil
			}
			continue
		}

		if checker == nil {
			checker = &apiRouteInfoChecker{
				start: i,
				lines: []string{line},
			}
			continue
		}

		checker.lines = append(checker.lines, line)
	}
	if checker != nil {
		checkers = append(checkers, checker)
	}

	syntaxErrors := []SyntaxError{}
	for _, checker := range checkers {
		checkErrors := checker.check()
		for _, checkError := range checkErrors {
			syntaxErrors = append(syntaxErrors, SyntaxError{
				Range: Range{
					Start: Position{
						Line:      uint32(checkError.line + checker.start),
						Character: uint32(checkError.start),
					},
					End: Position{
						Line:      uint32(checkError.line + checker.start),
						Character: uint32(checkError.end),
					},
				},
				Message: checkError.message,
			})
		}
	}

	return syntaxErrors
}

type apiRouteInfoChecker struct {
	start   int
	lines   []string
	tagDefs []tag.SwagTagDef
}

func (c apiRouteInfoChecker) hasSwagTags() bool {
	return len(c.tagDefs) > 0
}

func (c apiRouteInfoChecker) hasRouterTag() bool {
	for _, tagDef := range c.tagDefs {
		if tagDef.Type.IsRouter() {
			return true
		}
	}

	return false
}

type routeInfo struct {
	httpMethod string
	path       string
}

type checkError struct {
	message string
	line    int
	start   int
	end     int
}

func (c *apiRouteInfoChecker) check() []checkError {
	checkErrors := []checkError{}

	for lineIndex, line := range c.lines {
		firstToken, tokenizeArgs := parser.Tokenize(line)
		if !strings.HasPrefix(firstToken.Text, "@") {
			continue
		}
		if tokenizeArgs == nil {
			continue
		}
		swagTagDef := tag.NewSwagTagDef(strings.TrimPrefix(firstToken.Text, "@"))
		if !swagTagDef.IsValidTag() {
			continue
		}

		c.tagDefs = append(c.tagDefs, swagTagDef)

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
					line:    lineIndex,
					start:   argToken.Start,
					end:     argToken.End,
				})
			}

			i++
		}

		if i < swagTagDef.RequiredArgsCount() {
			checkErrors = append(checkErrors, checkError{
				message: swagTagDef.ErrorMessage(),
				line:    lineIndex,
				start:   firstToken.Start,
				end:     firstToken.End,
			})
		}
	}

	if c.hasSwagTags() && !c.hasRouterTag() {
		checkErrors = append(checkErrors, checkError{
			message: "@Router is required.",
			line:    0,
			start:   0,
			end:     0,
		})
	}

	return checkErrors
}
