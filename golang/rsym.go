package golang

import (
	"fmt"
	"github.com/golangee/src/ast"
	"github.com/golangee/src/render"
)

// renderSym emits an imported qualifier.
func (r *Renderer) renderSym(node *ast.Sym, w *render.BufferedWriter) error {
	switch node.Kind {
	case ast.SymTermStmt:
		w.Print(";")
	case ast.SymNewline:
		w.Print("\n")
	default:
		panic("unknown sym: " + fmt.Sprint(node.Kind))
	}

	return nil
}
