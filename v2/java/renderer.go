package java

import (
	"fmt"
	"github.com/golangee/src/v2"
	"path/filepath"
	"reflect"
	"strconv"
)

const mimeTypeJava = "application/java"
const packageJavaDocFile = "package-info"

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
	w.Printf("\n\n") // double line break, otherwise the formatter will purge it
	writeComment(w, file.parent.pkg.PackageName(), file.file.Doc())

	w.Printf("package %s;\n", file.parent.pkg.PackageName())
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
	w.Printf(visibilityAsKeyword(obj.Visibility()))
	w.Printf(" ")
	w.Printf("class %s {\n", obj.Name())
	w.Printf("}\n")

	return nil
}

func visibilityAsKeyword(v src.Visibility) string {
	switch v {
	case src.Public:
		return "public"
	case src.PackagePrivate:
		return ""
	case src.Private:
		return "private"
	case src.Protected:
		return "protected"
	default:
		panic("visibility not implemented: " + strconv.Itoa(int(v)))
	}
}

// Render emits the declared module as a java project.
func Render(mod *src.Module) ([]src.RenderedFile, error) {
	var res []src.RenderedFile

	tree := newModNode(mod)

	for _, p := range tree.packages {
		if p.pkg.Doc() != "" {
			pDoc := src.NewSrcFile(packageJavaDocFile)
			pDoc.SetDoc(p.pkg.Doc())
			pDoc.SetDocPreamble(p.pkg.DocPreamble())
			docNode := newSrcFileNode(p, pDoc)
			buf, err := renderFile(docNode)
			if err != nil {
				panic("illegal state: " + err.Error() + ": " + string(buf))
			}

			res = append(res, src.RenderedFile{
				AbsoluteFileName: filepath.Join(p.pkg.ImportPath(), docNode.file.Name()+".java"),
				MimeType:         mimeTypeJava,
				Buf:              buf,
				Error:            err,
			})

		}

		for _, file := range p.srcFiles {
			buf, err := renderFile(file)
			res = append(res, src.RenderedFile{
				AbsoluteFileName: filepath.Join(p.pkg.ImportPath(), file.file.Name()+".java"),
				MimeType:         mimeTypeJava,
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
