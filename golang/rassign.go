package golang

import (
	"fmt"
	"github.com/golangee/src/ast"
	"github.com/golangee/src/render"
	"strconv"
)

func (r *Renderer) renderAssignComment(node *ast.Assign, w *render.BufferedWriter) {
	if node.ObjComment != nil {
		var ellipsisName string
		if len(node.Lhs) > 0 {
			if ident, ok := node.Lhs[0].(*ast.Ident); ok {
				ellipsisName = ident.Name
			}
		}

		w.Print(formatComment(ellipsisName, node.ObjComment.Text))
	}
}

// renderAssign emits an assignment. If its actually an expression is language dependent.
func (r *Renderer) renderAssign(node *ast.Assign, w *render.BufferedWriter) error {
	if _, isConst := node.Parent().(*ast.ConstDecl); !isConst {
		r.renderAssignComment(node, w)
	}

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
