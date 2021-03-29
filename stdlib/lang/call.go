package lang

import "github.com/golangee/src/ast"

// CallStatic interprets the name as qualified and causes an import of the qualifier.
func CallStatic(name ast.Name, args ...ast.Expr) *ast.Macro {
	return ast.NewMacro().SetMatchers(
		ast.MatchTargetLanguage(ast.LangGo, ast.NewCallExpr(ast.NewSelExpr(ast.NewQualIdent(name.Qualifier()), ast.NewIdent(name.Identifier())), args...)),
	)
}

// CallIdent is like CallStatic but does not cause an import because it just uses local identifiers for the receiver and method.
func CallIdent(ident, method string, args ...ast.Expr) *ast.Macro {
	return ast.NewMacro().SetMatchers(
		ast.MatchTargetLanguage(ast.LangGo, ast.NewCallExpr(ast.NewSelExpr(ast.NewIdent(ident), ast.NewIdent(method)), args...)),
	)
}

// Call is like CallIdent but does not cause an import because it just uses a local identifier (like a static method).
func Call(ident string, args ...ast.Expr) *ast.Macro {
	return ast.NewMacro().SetMatchers(
		ast.MatchTargetLanguage(ast.LangGo, ast.NewCallExpr(ast.NewIdent(ident), args...)),
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
