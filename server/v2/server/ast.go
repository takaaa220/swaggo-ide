package server

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

func isInFunctionComment(pos token.Position) (bool, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, pos.Filename, nil, parser.ParseComments)
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
				if commentPos.Line == pos.Line {
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
