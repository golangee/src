package golang

import (
	"github.com/golangee/src/ast"
	"github.com/golangee/src/render"
)

func (r *Renderer) renderParam(node *ast.Param, w *render.BufferedWriter) error {
	w.Print(node.ParamName)
	w.Print(" ")
	
	return r.renderNode(node.ParamTypeDecl, w)
}
