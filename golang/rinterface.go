package golang

import (
	"fmt"
	"github.com/golangee/src/ast"
	"github.com/golangee/src/golang/validate"
	"github.com/golangee/src/render"
)

// renderStruct emits a struct type.
func (r *Renderer) renderInterface(node *ast.Interface, w *render.BufferedWriter) error {
	r.writeCommentNode(w, false, node.Identifier(), node.Comment())

	if _, ok := node.Parent().(*ast.File); ok {
		if err := validate.ExportedIdentifier(node.Visibility(), node.Identifier()); err != nil {
			return err
		}

		w.Printf(" type %s interface {\n", node.Identifier())
	} else {
		w.Printf(" %s interface {\n", node.Identifier())
	}

	/*
		for _, typeNode := range node.Types() {
			if err := r.renderType(typeNode, w); err != nil {
				return err
			}
		}*/

	for _, fun := range node.Methods() {
		funComment := r.renderFuncComment(fun)
		if funComment != "" {
			r.writeComment(w, false, fun.Identifier(), funComment)
		}

		if err := r.renderFunc(fun, w); err != nil {
			return fmt.Errorf("cannot render func '%s': %w", fun.Identifier(), err)
		}

		// I like a new line after a func but be more compact without comment
		if funComment != ""{
			w.Printf("\n")
		}
	}

	w.Printf("}\n")

	return nil
}
