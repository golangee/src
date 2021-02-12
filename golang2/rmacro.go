package golang

import (
	"fmt"
	"github.com/golangee/src/ast"
	"github.com/golangee/src/render"
)

// renderMacro emits something which is usually evaluated here.
func (r *Renderer) renderMacro(node *ast.Macro, w *render.BufferedWriter) error {
	r.writeCommentNode(w, false, "", node.Comment())
	if node.Func != nil {
		actualNodes := node.Func(node)
		for _, actualNode := range actualNodes {
			if err := r.renderNode(actualNode, w); err != nil {
				return fmt.Errorf("unable to render dynamic macro node: %w", err)
			}
		}
	}

	return nil
}
