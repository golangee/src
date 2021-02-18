package lang

import "github.com/golangee/src/ast"

func CallDefine(lhs, rhs ast.Expr) *ast.Macro {
	return ast.NewMacro().SetMatchers(
		ast.MatchTargetLanguage(ast.LangGo, ast.NewAssign(ast.Exprs(lhs), ast.AssignDefine, ast.Exprs(rhs))),
	)
}

func TryDefine(lhs, rhs ast.Expr, errMsg string) *ast.Macro {
	return ast.NewMacro().SetMatchers(
		ast.MatchTargetLanguage(ast.LangGo,
			ast.NewAssign(ast.Exprs(lhs, ast.NewIdent("err")), ast.AssignDefine, ast.Exprs(rhs)),
			Term(),


		),
	)
}
