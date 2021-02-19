package golang

import (
	"fmt"
	"github.com/golangee/src/ast"
	"github.com/golangee/src/golang/validate"
	"github.com/golangee/src/render"
)

// renderStruct emits a struct type.
func (r *Renderer) renderStruct(node *ast.Struct, w *render.BufferedWriter) error {
	r.writeCommentNode(w, false, node.Identifier(), node.Comment())

	if err := validate.ExportedIdentifier(node.Visibility(), node.Identifier()); err != nil {
		return err
	}

	w.Printf(" type %s struct {\n", node.Identifier())

	/*
		for _, typeNode := range node.Types() {
			if err := r.renderType(typeNode, w); err != nil {
				return err
			}
		}*/

	for _, field := range node.Fields() {
		if err := r.renderField(field, w); err != nil {
			return fmt.Errorf("cannot render field '%s': %w", field.Identifier(), err)
		}

		// unsure if want to keep a newline, but I find it more readable at least with comment
		if field.ObjComment != nil {
			w.Printf("\n")
		}
	}

	w.Printf("}\n")

	for _, fun := range node.Methods() {
		funComment := r.renderFuncComment(fun)
		if funComment != "" {
			r.writeComment(w, false, fun.Identifier(), funComment)
		}

		w.Printf("func ")
		if err := r.renderFunc(fun, w); err != nil {
			return fmt.Errorf("cannot render func '%s': %w", fun.Identifier(), err)
		}

		// I like a new line after a func
		w.Printf("\n")
	}

	return nil
}
