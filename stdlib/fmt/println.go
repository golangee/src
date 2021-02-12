package fmt

import "github.com/golangee/src/ast"

func Println(node ...ast.Node) *ast.Macro {
	return ast.NewMacro().SetMatchers(
		ast.MatchTargetLanguage(ast.LangGo, ast.NewSimpleTypeDecl("fmt.Println")),
	)
}
