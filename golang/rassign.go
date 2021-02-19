package golang

import (
	"fmt"
	"github.com/golangee/src/ast"
	"github.com/golangee/src/render"
	"strconv"
)

// renderAssign emits an assignment. If its actually an expression is language dependent.
func (r *Renderer) renderAssign(node *ast.Assign, w *render.BufferedWriter) error {
	for i, lh := range node.Lhs {
		if err := r.renderNode(lh, w); err != nil {
			return fmt.Errorf("unable to render lhs: %w", err)
		}

		if i < len(node.Lhs)-1 {
			w.Printf(", ")
		}
	}

	switch node.Kind {
	case ast.AssignSimple:
		w.Print("=")
	case ast.AssignDefine:
		w.Print(":=")
	case ast.AssignAdd:
		w.Print("+=")
	case ast.AssignSub:
		w.Print("-=")
	case ast.AssignMul:
		w.Print("*=")
	case ast.AssignRem:
		w.Print("%=")
	default:
		panic("assignment not implemented: " + strconv.Itoa(int(node.Kind)))
	}

	for i, rh := range node.Rhs {
		if err := r.renderNode(rh, w); err != nil {
			return fmt.Errorf("unable to render rhs: %w", err)
		}

		if i < len(node.Rhs)-1 {
			w.Printf(", ")
		}
	}

	return nil
}
