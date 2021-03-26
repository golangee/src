package golang

import (
	"fmt"
	"github.com/golangee/src/ast"
	"github.com/golangee/src/render"
)

// renderBinaryExpr emits a binary expression
func (r *Renderer) renderUnaryExpr(node *ast.UnaryExpr, w *render.BufferedWriter) error {
	switch node.Op {
	case ast.OpAdd:
		w.Print("+")
	case ast.OpSub:
		w.Print("-")
	case ast.OpAnd:
		w.Print("&")
	case ast.OpNot:
		w.Print("!")
	default:
		panic("operator not supported: " + fmt.Sprint(node.Op))
	}

	if err := r.renderNode(node.X, w); err != nil {
		return fmt.Errorf("unable to render x: %w", err)
	}

	return nil
}
