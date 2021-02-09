package golang

import (
	"github.com/golangee/src"
	"github.com/golangee/src/ast"
	"github.com/golangee/src/render"
	"strconv"
)

// renderFile generates the code for the entire file.
func (r *Renderer) renderFile(file *ast.File) ([]byte, error) {
	w := &render.BufferedWriter{}

	// package license or whatever
	if file.Pkg().Preamble != nil {
		r.writeComment(w, false, file.Pkg().Name, file.Pkg().Preamble.Text)
		w.Printf("\n\n") // double line break, otherwise the formatter will purge it
	}

	// file license or whatever
	if file.Preamble != nil {
		r.writeComment(w, false, file.Pkg().Name, file.Preamble.Text)
		w.Printf("\n\n") // double line break, otherwise the formatter will purge it
	}

	if file.Comment() != nil {
		r.writeComment(w, true, file.Pkg().Name, file.Comment().Text)
	}

	w.Printf("package %s;\n", file.Pkg().Name)

	// render everything into tmp first, the importer beautifies all required imports on-the-go
	tmp := &src.BufferedWriter{}
	for _, typ := range file.Types {
		if err := renderType(typ, tmp); err != nil {
			return nil, err
		}
	}

	for _, node := range file.Functions {
		if err := renderFunc(node, tmp); err != nil {
			return nil, err
		}
	}

	importer := r.importer(file)
	for namedImport, qualifier := range importer.namedImports {
		w.Printf("import %s %s\n", namedImport, strconv.Quote(qualifier))
	}

	w.Printf(tmp.String())

	return Format(w.Bytes())
}
