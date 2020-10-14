package java

import (
	"fmt"
	"github.com/golangee/src/v2"
	"github.com/golangee/src/v2/ast"
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
func renderFile(file *ast.SrcFileNode) ([]byte, error) {
	w := &src.BufferedWriter{}

	writeComment(w, file.PkgNode().SrcPackage().PackageName(), file.SrcFile().DocPreamble())
	w.Printf("\n\n") // double line break, otherwise the formatter will purge it
	writeComment(w, file.PkgNode().SrcPackage().PackageName(), file.SrcFile().Doc())

	w.Printf("package %s;\n", file.PkgNode().SrcPackage().PackageName())

	// render everything into tmp first, the importer beautifies all required imports on-the-go
	tmp := &src.BufferedWriter{}
	for _, typ := range file.Types() {
		if err := renderType(typ, tmp); err != nil {
			return nil, err
		}
	}

	importer := importerFromTree(file)
	for _, qualifier := range importer.qualifiers() {
		w.Printf("import %s;\n", qualifier)
	}

	w.Printf(tmp.String())

	return Format(w.Bytes())
}

func renderType(t *ast.TypeNode, w *src.BufferedWriter) error {
	switch node := t.NamedNode().(type) {
	case *ast.StructNode:
		return renderStruct(node, w)
	case *ast.InterfaceNode:
		return renderInterface(node, w)
	default:
		panic("type not yet implemented: " + reflect.TypeOf(t).String())
	}
}

func renderInterface(node *ast.InterfaceNode, w *src.BufferedWriter) error {
	writeComment(w, node.SrcInterface().Name(), node.SrcInterface().Doc())
	w.Printf(visibilityAsKeyword(node.SrcInterface().Visibility()))

	w.Printf(" interface %s {\n", node.SrcInterface().Name())
	for _, fun := range node.Methods() {
		if err := renderFunc(fun, w); err != nil {
			return fmt.Errorf("failed to render func %s: %w", fun.SrcFunc().Name(), err)
		}
	}
	w.Printf("}\n")

	return nil
}

func renderStruct(node *ast.StructNode, w *src.BufferedWriter) error {
	writeComment(w, node.SrcStruct().Name(), node.SrcStruct().Doc())
	w.Printf(visibilityAsKeyword(node.SrcStruct().Visibility()))
	if node.SrcStruct().Final() {
		w.Printf(" final ")
	}

	if node.SrcStruct().Static() {
		w.Printf(" static ")
	}

	w.Printf(" class %s {\n", node.SrcStruct().Name())
	for _, field := range node.Fields() {
		if err := renderField(field, w); err != nil {
			return fmt.Errorf("failed to render field %s: %w", field.SrcField().Name(), err)
		}
	}
	w.Printf("}\n")

	return nil
}

func renderFunc(node *ast.FuncNode, w *src.BufferedWriter) error {
	writeComment(w, node.SrcFunc().Name(), node.SrcFunc().Doc())
	w.Printf(visibilityAsKeyword(node.SrcFunc().Visibility()))
	w.Printf(" ")
	if len(node.SrcFunc().Results()) == 0 {
		w.Printf("void ")
	} else {
		/*(node * typeDeclNode, w * src.BufferedWriter)
		node.srcFunc.Results()[0].
			renderTypeDecl()*/
	}
	w.Printf(node.SrcFunc().Name())
	w.Printf("(")
	w.Printf(")")

	if node.SrcFunc().Body() == nil {
		w.Printf(";")
	} else {
		w.Printf("{\n")
		w.Printf("}\n")
	}

	return nil
}

func renderField(node *ast.FieldNode, w *src.BufferedWriter) error {
	writeComment(w, node.SrcField().Name(), node.SrcField().Doc())
	for _, annotation := range node.Annotations() {
		if err := renderAnnotation(annotation, w); err != nil {
			return err
		}
	}
	w.Printf(visibilityAsKeyword(node.SrcField().Visibility()))
	w.Printf(" ")
	if err := renderTypeDecl(node.TypeDecl(), w); err != nil {
		return err
	}
	w.Printf(" ")
	w.Printf(node.SrcField().Name())
	w.Printf(";\n")

	return nil
}

func renderAnnotation(node *ast.AnnotationNode, w *src.BufferedWriter) error {
	importer := importerFromTree(node)

	w.Printf("@")
	w.Printf(string(importer.shortify(node.SrcAnnotation().Name())))
	attrs := node.SrcAnnotation().Attributes()
	if len(attrs) > 0 {
		w.Printf("(")
		// the default case
		if len(attrs) == 1 && attrs[0] == "" {
			w.Printf(node.SrcAnnotation().Value(""))
		} else {
			// the named attribute cases
			for i, attr := range attrs {
				w.Printf(attr)
				w.Printf(" = ")
				w.Printf(node.SrcAnnotation().Value(attr))
				if i < len(attrs)-1 {
					w.Printf(", ")
				}
			}
		}

		w.Printf(")")
	}

	return nil
}

func renderTypeDecl(node ast.TypeDeclNode, w *src.BufferedWriter) error {
	importer := importerFromTree(node)

	switch t := node.(type) {
	case *ast.SimpleTypeDeclNode:
		w.Printf(string(importer.shortify(fromStdlib(t.SrcSimpleTypeDecl().Name()))))
	case *ast.TypeDeclPtrNode:
		atomicReference := importer.shortify("java.util.concurrent.atomic.AtomicReference")
		w.Printf(string(atomicReference) + "<")
		if err := renderTypeDecl(t.TypeDecl(), w); err != nil {
			return err
		}
		w.Printf(">")
	case *ast.SliceTypeDeclNode:
		if err := renderTypeDecl(t.TypeDecl(), w); err != nil {
			return err
		}
		w.Printf("[]")
	case *ast.GenericTypeDeclNode:
		if err := renderTypeDecl(t.TypeDecl(), w); err != nil {
			return err
		}
		w.Printf("<")
		for i, decl := range t.Params() {
			if err := renderTypeDecl(decl, w); err != nil {
				return err
			}
			if i < len(t.Params())-1 {
				w.Printf(",")
			}
		}
		w.Printf(">")
	case *ast.ChanTypeDeclNode:
		blockingQueue := importer.shortify("java.util.concurrent.BlockingQueue")
		w.Printf(string(blockingQueue) + "<")
		if err := renderTypeDecl(t.TypeDecl(), w); err != nil {
			return err
		}
		w.Printf(">")

	case *ast.ArrayTypeDeclNode:
		// in Java this is the same as a slice, we cannot have yet custom size value arrays. Perhaps
		// valhalla may fix that
		if err := renderTypeDecl(t.TypeDecl(), w); err != nil {
			return err
		}
		w.Printf("[]")
	case *ast.FuncTypeDeclNode:
		// Java does not have it. We would need to create a functional interface for it, which is out of scope here.
		writeComment(w, "", "inline function declarations are not supported by java:\n\n"+t.SrcFuncTypeDecl().String())
		w.Printf("Object")
		return nil
	default:
		panic("not yet implemented: " + reflect.TypeOf(t).String())
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

	tree := ast.NewModNode(mod)
	installImporter(tree)

	for _, p := range tree.Packages() {
		if p.SrcPackage().Doc() != "" {
			pDoc := src.NewSrcFile(packageJavaDocFile)
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
				MimeType:         mimeTypeJava,
				Buf:              buf,
				Error:            err,
			})

		}

		for _, file := range p.Files() {
			buf, err := renderFile(file)
			res = append(res, src.RenderedFile{
				AbsoluteFileName: filepath.Join(p.SrcPackage().ImportPath(), file.SrcFile().Name()+".java"),
				MimeType:         mimeTypeJava,
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
