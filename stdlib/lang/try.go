package lang

import (
	"fmt"
	"github.com/golangee/src/ast"
	"github.com/golangee/src/stdlib"
)

// CallDefine emits a variable (re)declaration with an assignment.
func CallDefine(lhs, rhs ast.Expr) *ast.Macro {
	return ast.NewMacro().SetMatchers(
		ast.MatchTargetLanguage(ast.LangGo, ast.NewAssign(ast.Exprs(lhs), ast.AssignDefine, ast.Exprs(rhs))),
	)
}

// TryDefine emits a variable (re)declaration with an assignment and an error check with early return.
// It evaluates the current context to decide how to return and how to re-throw error.
func TryDefine(lhs, rhs ast.Expr, errMsg string) *ast.Macro {
	return ast.NewMacro().SetMatchers(
		ast.MatchTargetLanguageWithContext(ast.LangGo,
			func(m *ast.Macro) []ast.Node {
				myFunc := assertFunc(m)
				if len(myFunc.FunResults) == 0 {
					panic("func " + myFunc.FunName + " must define at least an error return value")
				}

				lastSTD, ok := myFunc.FunResults[len(myFunc.FunResults)-1].ParamTypeDecl.(*ast.SimpleTypeDecl)
				if !ok || lastSTD.SimpleName != stdlib.Error {
					panic("func " + myFunc.FunName + " last result must be an error return value but is " + fmt.Sprint(lastSTD))
				}

				var results []ast.Expr
				for i := 0; i < len(myFunc.FunResults)-1; i++ {
					decl := myFunc.FunResults[i].TypeDecl()
					switch t := decl.(type) {
					case *ast.SimpleTypeDecl:
						if stdlib.IsNumber(string(t.SimpleName)) {
							results = append(results, ast.NewIntLit(0))
						} else if t.SimpleName == stdlib.String {
							results = append(results, ast.NewStrLit(""))
						} else
						{
							// TODO this is not always correct and cannot always be resolved due to external dependencies
							results = append(results, ast.NewCompLit(t))
						}
					case *ast.SliceTypeDecl:
						results = append(results, ast.NewIdent("nil"))
					case *ast.TypeDeclPtr:
						results = append(results, ast.NewIdent("nil"))
					}

				}

				results = append(results, CallStatic("fmt.Errorf", ast.NewStrLit(errMsg+": %w"), ast.NewIdent("err")))

				return ast.Nodes(
					ast.NewAssign(ast.Exprs(lhs, ast.NewIdent("err")), ast.AssignDefine, ast.Exprs(rhs)),
					Term(),
					ast.NewIfStmt(ast.NewBinaryExpr(ast.NewIdent("err"), ast.OpNotEqual, ast.NewIdent("nil")), ast.NewBlock(
						ast.NewReturnStmt(results...),
					)),
				)
			},

		),
	)
}

// there is always an outer func definition
func assertFunc(n ast.Node) *ast.Func {
	f := &ast.Func{}
	if ok := ast.ParentAs(n, &f); ok {
		return f
	}

	panic("invalid context: must be a func child")
}
