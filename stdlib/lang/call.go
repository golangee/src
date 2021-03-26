package lang

import "github.com/golangee/src/ast"

// CallStatic interprets the name as qualified and causes an import of the qualifier.
func CallStatic(name ast.Name, args ...ast.Expr) *ast.Macro {
	return ast.NewMacro().SetMatchers(
		ast.MatchTargetLanguage(ast.LangGo, ast.NewCallExpr(ast.NewSelExpr(ast.NewQualIdent(name.Qualifier()), ast.NewIdent(name.Identifier())), args...)),
	)
}

// CreateLiteral takes the
func CreateLiteral(name ast.Name, args ...ast.Expr) *ast.Macro {
	return ast.NewMacro().SetMatchers(
		ast.MatchTargetLanguage(ast.LangGo, ast.NewCompLit(ast.NewSimpleTypeDecl(name), args...)),
	)
}

// ToString converts the given expression into a string.
func ToString(expr ast.Expr) *ast.Macro {
	return ast.NewMacro().SetMatchers(
		ast.MatchTargetLanguage(ast.LangGo, CallStatic("fmt.Sprintf", ast.NewStrLit("%v"), expr)),
	)
}

// Itoa performs a more optimized integer to ascii.
func Itoa(expr ast.Expr) *ast.Macro {
	return ast.NewMacro().SetMatchers(
		ast.MatchTargetLanguage(ast.LangGo, CallStatic("strconv.Itoa", expr)),
	)
}
