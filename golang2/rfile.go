package golang

import (
	"github.com/golangee/src/ast"
	"github.com/golangee/src/render"
	"strconv"
)

// renderFile generates the code for the entire file.
func (r *Renderer) renderFile(file *ast.File) ([]byte, error) {
	w := &render.BufferedWriter{}

	// file license or whatever
	if file.Preamble != nil {
		r.writeComment(w, false, file.Pkg().Name, file.Preamble.Text)
		w.Printf("\n\n") // double line break, otherwise the formatter will purge it
	}

	// actual package comment
	if file.Comment() != nil {
		r.writeComment(w, true, file.Pkg().Name, file.Comment().Text)
	}

	w.Printf("package %s\n", file.Pkg().Name)

	// render everything into tmp first, the importer beautifies all required imports on-the-go
	tmp := &render.BufferedWriter{}
	for _, typ := range file.Types {
		if err := r.renderNode(typ, tmp); err != nil {
			return nil, err
		}
	}

	for _, fun := range file.Functions {
		funComment := r.renderFuncComment(fun)
		if funComment != "" {
			r.writeComment(tmp, false, fun.Identifier(), funComment)
		}

		if err := r.renderFunc(fun, tmp); err != nil {
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
