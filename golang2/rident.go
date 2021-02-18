package golang

import (
	"github.com/golangee/src/ast"
	"github.com/golangee/src/render"
)

// renderIdent emits an identifiers name.
func (r *Renderer) renderIdent(node *ast.Ident, w *render.BufferedWriter) error {
	w.Printf(node.Name)

	return nil
}
