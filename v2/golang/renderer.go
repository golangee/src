package golang

import (
	"fmt"
	"github.com/golangee/src/v2"
	"path/filepath"
	"reflect"
)

const packageGoDocFile = "doc"
const mimeTypeGo = "application/go"

func writeComment(w *src.BufferedWriter, name, doc string) {
	myDoc := formatComment(name, doc)
	if doc != "" {
		w.Printf(myDoc)
	}
}

// RenderFile tries to emit the file as java
func renderFile(file *srcFileNode) ([]byte, error) {
	w := &src.BufferedWriter{}

	writeComment(w, file.parent.pkg.PackageName(), file.file.DocPreamble())
	w.Printf("\n")
	writeComment(w, file.parent.pkg.PackageName(), file.file.Doc())

	w.Printf("package %s\n", file.parent.pkg.PackageName())
	for _, typ := range file.types {
		if err := renderType(typ, w); err != nil {
			return nil, err
		}
	}
	return Format(w.Bytes())
}

func renderType(t *typeNode, w *src.BufferedWriter) error {
	switch obj := t.namedType.(type) {
	case *src.Struct:
		return renderStruct(t, obj, w)
	default:
		panic("type not yet implemented: " + reflect.TypeOf(t).String())
	}
}

func renderStruct(node *typeNode, obj *src.Struct, w *src.BufferedWriter) error {
	writeComment(w,obj.Name(),obj.Doc())
	w.Printf("type ")
	if obj.Visibility() == src.Public {
		w.Printf(MakePublic(obj.Name()))
	} else {
		w.Printf(MakePrivate(obj.Name()))
	}
	w.Printf(" ")
	w.Printf("struct {\n")
	w.Printf("}\n")

	return nil
}

// Render emits the declared module as a java project.
func Render(mod *src.Module) ([]src.RenderedFile, error) {
	var res []src.RenderedFile

	tree := newModNode(mod)

	for _, p := range tree.packages {
		if p.pkg.Doc() != "" {
			pDoc := src.NewSrcFile(packageGoDocFile)
			pDoc.SetDoc(p.pkg.Doc())
			pDoc.SetDocPreamble(p.pkg.DocPreamble())
			docNode := newSrcFileNode(p, pDoc)
			buf, err := renderFile(docNode)
			if err != nil {
				panic("illegal state: " + err.Error() + ": " + string(buf))
			}

			res = append(res, src.RenderedFile{
				AbsoluteFileName: filepath.Join(p.pkg.ImportPath(), docNode.file.Name()+".go"),
				MimeType:         mimeTypeGo,
				Buf:              buf,
				Error:            err,
			})

		}

		for _, file := range p.srcFiles {
			buf, err := renderFile(file)
			res = append(res, src.RenderedFile{
				AbsoluteFileName: filepath.Join(p.pkg.ImportPath(), file.file.Name()+".go"),
				MimeType:         mimeTypeGo,
				Buf:              buf,
				Error:            err,
			})

			if err != nil {
				return res, fmt.Errorf("unable to render %s/%s: %w", p.pkg.ImportPath(), file.file.Name(), err)
			}
		}
	}

	return res, nil
}
