package server

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/takaaa220/go-swag-ide/server/v2/server-sdk/protocol"
)

func isInFunctionComment(filePath string, src string, pos protocol.Position) (bool, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, src, parser.ParseComments)
	if err != nil {
		return false, err
	}

	for _, decl := range node.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if funcDecl.Doc == nil {
				continue
			}

			for _, comment := range funcDecl.Doc.List {
				commentPos := fset.Position(comment.Pos())
				if commentPos.Line == int(pos.Line+1) {
					commentText := strings.TrimSpace(comment.Text)
					if strings.HasPrefix(commentText, "//") {
						return true, nil
					}
				}
			}
		}
	}

	return false, nil
}
