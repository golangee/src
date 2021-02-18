package lang

import "github.com/golangee/src/ast"

func CallStatic(name ast.Name, args ...ast.Expr) *ast.Macro {
	return ast.NewMacro().SetMatchers(
		ast.MatchTargetLanguage(ast.LangGo, ast.NewCallExpr(ast.NewSelExpr(ast.NewQualIdent(name.Qualifier()), ast.NewIdent(name.Identifier())), args...)),
	)
}
