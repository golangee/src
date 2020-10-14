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
		w.Printf("\n")
	}
}

// RenderFile tries to emit the file as java
func renderFile(file *srcFileNode) ([]byte, error) {
	w := &src.BufferedWriter{}

	writeComment(w, file.parent.srcPackage.PackageName(), file.srcFile.DocPreamble())
	w.Printf("\n\n") // double line break, otherwise the formatter will purge it
	writeComment(w, file.parent.srcPackage.PackageName(), file.srcFile.Doc())

	w.Printf("package %s;\n", file.parent.srcPackage.PackageName())

	// render everything into tmp first, the importer beautifies all required imports on-the-go
	tmp := &src.BufferedWriter{}
	for _, typ := range file.types {
		if err := renderType(typ, tmp); err != nil {
			return nil, err
		}
	}

	for _, qualifier := range file.importer.qualifiers() {
		w.Printf("import %s;\n", qualifier)
	}

	w.Printf(tmp.String())

	return Format(w.Bytes())
}

func renderType(t *typeNode, w *src.BufferedWriter) error {
	switch node := t.namedNode.(type) {
	case *structNode:
		return renderStruct(node, w)
	default:
		panic("type not yet implemented: " + reflect.TypeOf(t).String())
	}
}

func renderStruct(node *structNode, w *src.BufferedWriter) error {
	writeComment(w, node.srcStruct.Name(), node.srcStruct.Doc())
	w.Printf(visibilityAsKeyword(node.srcStruct.Visibility()))
	w.Printf(" ")
	w.Printf("class %s {\n", node.srcStruct.Name())
	for _, field := range node.fields {
		if err := renderField(field, w); err != nil {
			return fmt.Errorf("failed to render field %s: %w", field.srcField.Name(), err)
		}
	}
	w.Printf("}\n")

	return nil
}

func renderField(node *fieldNode, w *src.BufferedWriter) error {
	writeComment(w, node.srcField.Name(), node.srcField.Doc())
	for _, annotation := range node.annotations {
		if err := renderAnnotation(annotation, w); err != nil {
			return err
		}
	}
	w.Printf(visibilityAsKeyword(node.srcField.Visibility()))
	w.Printf(" ")
	if err := renderTypeDecl(node.typeDeclNode, w); err != nil {
		return err
	}
	w.Printf(" ")
	w.Printf(node.srcField.Name())
	w.Printf(";\n")

	return nil
}

func renderAnnotation(node *annotationNode, w *src.BufferedWriter) error {
	importer := importerFromTree(node)

	w.Printf("@")
	w.Printf(string(importer.shortify(node.srcAnnotation.Name())))
	attrs := node.srcAnnotation.Attributes()
	if len(attrs) > 0 {
		w.Printf("(")
		// the default case
		if len(attrs) == 1 && attrs[0] == "" {
			w.Printf(node.srcAnnotation.Value(""))
		} else {
			// the named attribute cases
			for i, attr := range attrs {
				w.Printf(attr)
				w.Printf(" = ")
				w.Printf(node.srcAnnotation.Value(attr))
				if i < len(attrs)-1 {
					w.Printf(", ")
				}
			}
		}

		w.Printf(")")
	}

	return nil
}

func renderTypeDecl(node *typeDeclNode, w *src.BufferedWriter) error {
	importer := importerFromTree(node)

	switch t := node.srcTypeDecl.(type) {
	case *src.SimpleTypeDecl:
		w.Printf(string(importer.shortify(fromStdlib(t.Name()))))
	case *src.TypeDeclPtr:
		atomicReference := importer.shortify("java.util.concurrent.atomic.AtomicReference")
		w.Printf(string(atomicReference) + "<")
		childTypeDecl := newTypeDeclNode(node, t.TypeDecl())
		if err := renderTypeDecl(childTypeDecl, w); err != nil {
			return err
		}
		w.Printf(">")
	case *src.SliceTypeDecl:
		childTypeDecl := newTypeDeclNode(node, t.TypeDecl())
		if err := renderTypeDecl(childTypeDecl, w); err != nil {
			return err
		}
		w.Printf("[]")
	case *src.GenericTypeDecl:
		baseType := newTypeDeclNode(node, t.TypeDecl())
		if err := renderTypeDecl(baseType, w); err != nil {
			return err
		}
		w.Printf("<")
		for i, decl := range t.Params() {
			childTypeDecl := newTypeDeclNode(node, decl)
			if err := renderTypeDecl(childTypeDecl, w); err != nil {
				return err
			}
			if i < len(t.Params())-1 {
				w.Printf(",")
			}
		}
		w.Printf(">")
	case *src.ChanTypeDecl:
		blockingQueue := importer.shortify("java.util.concurrent.BlockingQueue")
		w.Printf(string(blockingQueue) + "<")
		childTypeDecl := newTypeDeclNode(node, t.TypeDecl())
		if err := renderTypeDecl(childTypeDecl, w); err != nil {
			return err
		}
		w.Printf(">")

	case *src.ArrayTypeDecl:
		// in java this is the same as a slice, we cannot have yet custom size value array. Perhaps
		// valhalla may fix that
		childTypeDecl := newTypeDeclNode(node, t.TypeDecl())
		if err := renderTypeDecl(childTypeDecl, w); err != nil {
			return err
		}
		w.Printf("[]")

	default:
		panic("not yet implemented: " + reflect.TypeOf(node.srcTypeDecl).String())
	}

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
		if p.srcPackage.Doc() != "" {
			pDoc := src.NewSrcFile(packageJavaDocFile)
			pDoc.SetDoc(p.srcPackage.Doc())
			pDoc.SetDocPreamble(p.srcPackage.DocPreamble())
			docNode := newSrcFileNode(p, pDoc)
			buf, err := renderFile(docNode)
			if err != nil {
				panic("illegal state: " + err.Error() + ": " + string(buf))
			}

			res = append(res, src.RenderedFile{
				AbsoluteFileName: filepath.Join(p.srcPackage.ImportPath(), docNode.srcFile.Name()+".java"),
				MimeType:         mimeTypeJava,
				Buf:              buf,
				Error:            err,
			})

		}

		for _, file := range p.files {
			buf, err := renderFile(file)
			res = append(res, src.RenderedFile{
				AbsoluteFileName: filepath.Join(p.srcPackage.ImportPath(), file.srcFile.Name()+".java"),
				MimeType:         mimeTypeJava,
				Buf:              buf,
				Error:            err,
			})

			if err != nil {
				return res, fmt.Errorf("unable to render %s/%s: %w", p.srcPackage.ImportPath(), file.srcFile.Name(), err)
			}
		}
	}

	return res, nil
}
