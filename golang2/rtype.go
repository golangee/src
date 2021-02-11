package golang

import (
	"github.com/golangee/src/ast"
	"github.com/golangee/src/render"
	"reflect"
)

// renderType inspects and emits the actual type.
func (r *Renderer) renderType(node ast.Node, w *render.BufferedWriter) error {
	switch n := node.(type) {
	case *ast.Struct:
		return r.renderStruct(n, w)
	case *ast.Func:
		return r.renderFunc(n, w)
	default:
		panic("unsupported type: " + reflect.TypeOf(n).String())
	}
}
