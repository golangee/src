package golang

import (
	"github.com/golangee/src/ast"
	"github.com/golangee/src/render"
)

func (r *Renderer) renderDeferStmt(node *ast.DeferStmt, w *render.BufferedWriter) error {
	w.Print("defer ")

	return r.renderNode(node.CallExpr, w)
}
