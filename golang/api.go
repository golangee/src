package golang

import (
	"fmt"
	"github.com/golangee/src"
	"github.com/golangee/src/ast"
	"path/filepath"
)

const packageGoDocFile = "doc.go"
const mimeTypeGo = "text/x-go-source"

// Render emits the declared module as a go project.
func Render(mod *src.Module) ([]src.RenderedFile, error) {
	var res []src.RenderedFile

	tree := ast.NewModNode(mod)
	installImporter(tree)

	for _, p := range tree.Packages() {
		if p.SrcPackage().Doc() != "" {
			pDoc := src.NewSrcFile(packageGoDocFile)
			pDoc.SetDoc(p.SrcPackage().Doc())
			pDoc.SetDocPreamble(p.SrcPackage().DocPreamble())
			docNode := ast.NewSrcFileNode(p, pDoc)
			docNode.SetValue(importerId, newImporter())
			buf, err := renderFile(docNode)
			if err != nil {
				panic("illegal state: " + err.Error() + ": " + string(buf))
			}

			res = append(res, src.RenderedFile{
				AbsoluteFileName: filepath.Join(p.SrcPackage().ImportPath(), docNode.SrcFile().Name()+".java"),
				MimeType:         mimeTypeGo,
				Buf:              buf,
				Error:            err,
			})

		}

		for _, file := range p.Files() {
			buf, err := renderFile(file)
			res = append(res, src.RenderedFile{
				AbsoluteFileName: filepath.Join(p.SrcPackage().ImportPath(), file.SrcFile().Name()+".java"),
				MimeType:         mimeTypeGo,
				Buf:              buf,
				Error:            err,
			})

			if err != nil {
				return res, fmt.Errorf("unable to render %s/%s: %w", p.SrcPackage().ImportPath(), file.SrcFile().Name(), err)
			}
		}
	}

	return res, nil
}
