package golang

import (
	"github.com/golangee/src/ast"
	"github.com/golangee/src/golang2/validate"
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

	/*
		for _, field := range node.Fields() {
			if err := renderField(field, w); err != nil {
				return fmt.Errorf("failed to render field %s: %w", field.SrcField().Name(), err)
			}
		}*/

	w.Printf("}\n")
	/*
		for _, fun := range node.Methods() {
			w.Printf("func ")
			if err := renderFunc(fun, w); err != nil {
				return fmt.Errorf("failed to render func %s: %w", fun.SrcFunc().Name(), err)
			}
		}*/

	return nil
}
