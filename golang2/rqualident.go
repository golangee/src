package golang

import (
	"github.com/golangee/src/ast"
	"github.com/golangee/src/render"
)

// renderQualIdent emits an imported qualifier.
func (r *Renderer) renderQualIdent(node *ast.QualIdent, w *render.BufferedWriter) error {
	importer := r.importer(node)
	renamedQualifier := importer.shortify(fromStdlib(ast.Name(node.Qualifier) + "._")).Qualifier()
	w.Printf(renamedQualifier)

	return nil
}
