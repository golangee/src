package golang

import (
	"fmt"
	"github.com/golangee/src/ast"
	"github.com/golangee/src/render"
	"reflect"
)

// renderType inspects and emits the actual type.
func (r *Renderer) renderNode(node ast.Node, w *render.BufferedWriter) error {
	switch n := node.(type) {
	case *ast.Struct:
		if err := r.renderStruct(n, w); err != nil {
			return fmt.Errorf("cannot render struct '%s': %w", n.Identifier(), err)
		}
	case *ast.Func:
		return r.renderFunc(n, w)
	case *ast.Block:
		return r.renderBlock(n, w)
	case *ast.Macro:
		if err := r.renderMacro(n, w); err != nil {
			return fmt.Errorf("cannot render macro: %w",  err)
		}
	case ast.TypeDecl:
		if err := r.renderTypeDecl(n, w); err != nil {
			return fmt.Errorf("cannot render TypeDecl: %w",  err)
		}
	default:
		panic("unsupported type: " + reflect.TypeOf(n).String())
	}

	return nil
}
