package fmt

import (
	"github.com/golangee/src/ast"
	"github.com/golangee/src/stdlib/lang"
)

func Println(args ...ast.Expr) *ast.Macro {
	return lang.CallStatic("fmt.Println", args...)
}
