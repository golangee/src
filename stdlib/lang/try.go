package lang

import "github.com/golangee/src/ast"

func CallDefine(lhs, rhs ast.Expr) *ast.Macro {
	return ast.NewMacro().SetMatchers(
		ast.MatchTargetLanguage(ast.LangGo, ast.NewAssign(ast.Exprs(lhs), ast.AssignDefine, ast.Exprs(rhs))),
	)
}

func TryDefine(lhs, rhs ast.Expr, errMsg string) *ast.Macro {
	// TODO inspect outer context to pick up correct behavior
	return ast.NewMacro().SetMatchers(
		ast.MatchTargetLanguage(ast.LangGo,
			ast.NewAssign(ast.Exprs(lhs, ast.NewIdent("err")), ast.AssignDefine, ast.Exprs(rhs)),
			Term(),
			ast.NewIfStmt(ast.NewBinaryExpr(ast.NewIdent("err"), ast.OpNotEqual, ast.NewIdent("nil")), ast.NewBlock(
				ast.NewReturnStmt(ast.NewIdent("nil"), CallStatic("fmt.Errorf", ast.NewStrLit(errMsg+": %w"), ast.NewIdent("err"))),
			)),

		),
	)
}
