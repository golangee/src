package golang

import (
	"fmt"
	"github.com/golangee/src/v2"
	"github.com/golangee/src/v2/ast"
	"github.com/golangee/src/v2/stdlib"
	"reflect"
	"strconv"
	"strings"
)

func writeComment(w *src.BufferedWriter, isPkg bool, name, doc string) {
	if isPkg {
		name = "Package " + name
	}

	myDoc := formatComment(name, doc)
	if doc != "" {
		w.Printf(myDoc)
	}
}

// renderFile tries to emit the file as java
func renderFile(file *ast.SrcFileNode) ([]byte, error) {
	w := &src.BufferedWriter{}

	writeComment(w, false, file.PkgNode().SrcPackage().PackageName(), file.SrcFile().DocPreamble())
	w.Printf("\n\n") // double line break, otherwise the formatter will purge it

	writeComment(w, true, file.PkgNode().SrcPackage().PackageName(), file.SrcFile().Doc())
	w.Printf("package %s;\n", file.PkgNode().SrcPackage().PackageName())

	// render everything into tmp first, the importer beautifies all required imports on-the-go
	tmp := &src.BufferedWriter{}
	for _, typ := range file.Types() {
		if err := renderType(typ, tmp); err != nil {
			return nil, err
		}
	}

	for _, node := range file.Functions() {
		if err := renderFunc(node, tmp); err != nil {
			return nil, err
		}
	}

	importer := importerFromTree(file)
	for namedImport, qualifier := range importer.namedImports {
		w.Printf("import %s %s\n", namedImport, strconv.Quote(qualifier))
	}

	w.Printf(tmp.String())

	return Format(w.Bytes())
}

func renderTypePreamble(w *src.BufferedWriter, node interface {
	Name() string
	Doc() string
	Annotations() []*ast.AnnotationNode
}) error {
	writeComment(w, false, node.Name(), node.Doc())

	for _, annotation := range node.Annotations() {
		if err := renderJavaLikeAnnotation(annotation, w); err != nil {
			return err
		}
		w.Printf("\n")
	}

	return nil
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
	if err := renderTypePreamble(w, node); err != nil {
		return err
	}

	w.Printf(" type %s interface{\n", visibilityAsName(node.SrcInterface().Visibility(), node.SrcInterface().Name()))

	for _, fun := range node.Methods() {
		if err := renderFunc(fun, w); err != nil {
			return fmt.Errorf("failed to render func %s: %w", fun.SrcFunc().Name(), err)
		}
	}
	w.Printf("}\n")

	// go does not have interface embedded type definition (besides anonymous types which are part of the function
	// declaration itself
	for _, typeNode := range node.Types() {
		if err := renderType(typeNode, w); err != nil {
			return err
		}
	}

	return nil
}

func renderStruct(node *ast.StructNode, w *src.BufferedWriter) error {
	if err := renderTypePreamble(w, node); err != nil {
		return err
	}

	w.Printf(" type %s struct {\n", visibilityAsName(node.SrcStruct().Visibility(), node.SrcStruct().Name()))

	for _, typeNode := range node.Types() {
		if err := renderType(typeNode, w); err != nil {
			return err
		}
	}

	for _, field := range node.Fields() {
		if err := renderField(field, w); err != nil {
			return fmt.Errorf("failed to render field %s: %w", field.SrcField().Name(), err)
		}
	}

	w.Printf("}\n")

	for _, fun := range node.Methods() {
		w.Printf("func ")
		if err := renderFunc(fun, w); err != nil {
			return fmt.Errorf("failed to render func %s: %w", fun.SrcFunc().Name(), err)
		}
	}

	return nil
}

func renderFunc(node *ast.FuncNode, w *src.BufferedWriter) error {
	comment := &strings.Builder{}
	comment.WriteString(node.SrcFunc().Doc())
	for _, annotation := range node.Annotations() {
		if err := renderJavaLikeAnnotation(annotation, w); err != nil {
			return err
		}
		w.Printf("\n")
	}

	for _, parameterNode := range node.InputParams() {
		for _, annotationNode := range parameterNode.Annotations() {
			if err := renderJavaLikeAnnotation(annotationNode, w); err != nil {
				return err
			}

			w.Printf("\n")
		}
	}
	comment.WriteString("\n")

	for _, parameterNode := range node.InputParams() {
		if parameterNode.SrcParameter().Doc() == "" {
			continue
		}

		comment.WriteString("@param ")
		name := parameterNode.SrcParameter().Name()
		if name == "" {
			name = fromStdlib(src.Name(parameterNode.SrcParameter().TypeDecl().String())).Identifier()
		}

		comment.WriteString(deEllipsis(name, parameterNode.SrcParameter().Doc()))
		comment.WriteString("\n")
	}

	for i, parameterNode := range node.OutputParams() {
		if i == 0 || parameterNode.SrcParameter().Doc() == "" {
			continue
		}

		comment.WriteString("@throws ")
		name := parameterNode.SrcParameter().Name()
		if name == "" {
			name = fromStdlib(src.Name(parameterNode.SrcParameter().TypeDecl().String())).Identifier()
		}

		comment.WriteString(deEllipsis(name, parameterNode.SrcParameter().Doc()))
		comment.WriteString("\n")
	}

	var structNode *ast.StructNode
	switch t := node.Parent().(type) {
	case *ast.StructNode:
		structNode = t
		recName := "_"
		if node.SrcFunc().RecName() != "" {
			recName = node.SrcFunc().RecName()
		}

		if node.SrcFunc().PtrReceiver() {
			recName += "*"
		}

		w.Printf("(%s %s) ", recName, t.Name())
	case *ast.SrcFileNode:
		w.Printf("func ")
	}

	w.Printf(visibilityAsName(node.SrcFunc().Visibility(), node.SrcFunc().Name()))
	w.Printf("(")
	for i, parameterNode := range node.InputParams() {

		if err := renderTypeDecl(parameterNode.TypeDecl(), w); err != nil {
			return err
		}

		if i == len(node.InputParams())-1 && node.SrcFunc().Variadic() {
			w.Printf("...")
		} else {
			w.Printf(" ")
		}

		w.Printf(parameterNode.SrcParameter().Name())

		if i < len(node.OutputParams())-1 {
			w.Printf(", ")
		}
	}
	w.Printf(")")

	if len(node.OutputParams()) > 0 {
		w.Printf("(")
	}
	for i, parameterNode := range node.OutputParams() {
		if err := renderTypeDecl(parameterNode.TypeDecl(), w); err != nil {
			return err
		}

		if i < len(node.OutputParams())-1 {
			w.Printf(", ")
		}
	}
	if len(node.OutputParams()) > 0 {
		w.Printf(")")
	}

	if node.SrcFunc().Body() == nil && structNode == nil {
		w.Printf("\n")
	} else {
		w.Printf("{\n")
		w.Printf("}\n")
	}

	return nil
}

func renderField(node *ast.FieldNode, w *src.BufferedWriter) error {
	writeComment(w, false, node.SrcField().Name(), node.SrcField().Doc())
	for _, annotation := range node.Annotations() {
		if err := renderJavaLikeAnnotation(annotation, w); err != nil {
			return err
		}
		w.Printf("\n")
	}

	w.Printf(node.SrcField().Name())
	w.Printf(" ")
	if err := renderTypeDecl(node.TypeDecl(), w); err != nil {
		return err
	}

	// try to translate annotations on fields into go struct tags
	if len(node.Annotations()) > 0 {
		w.Printf(" `")
		for i, annotationNode := range node.Annotations() {
			w.Printf(string(annotationNode.SrcAnnotation().Name()))
			w.Printf(":")
			// by definition the go renderer ever uses the empty value as is
			v := annotationNode.SrcAnnotation().Value("")
			w.Printf(strconv.Quote(v))

			if i < len(node.Annotations())-1 {
				w.Printf(" ") // separator between tag fields is the white space
			}
		}
		w.Printf("`")
	}

	w.Printf("\n")

	return nil
}

func renderJavaLikeAnnotation(node *ast.AnnotationNode, w *src.BufferedWriter) error {
	importer := importerFromTree(node)

	w.Printf("//@")
	if pn, ok := node.Parent().(*ast.ParameterNode); ok {
		w.Printf("[%s]", pn.SrcParameter().Name())
	}
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
		w.Printf("*")
		if err := renderTypeDecl(t.TypeDecl(), w); err != nil {
			return err
		}
	case *ast.SliceTypeDeclNode:
		w.Printf("[]")
		if err := renderTypeDecl(t.TypeDecl(), w); err != nil {
			return err
		}
	case *ast.GenericTypeDeclNode:
		if err := renderTypeDecl(t.TypeDecl(), w); err != nil {
			return err
		}

		builtInHandled := false
		if std, ok := t.SrcGenericTypeDecl().TypeDecl().(*src.SimpleTypeDecl); ok {
			switch std.Name() {
			case stdlib.Map:
				w.Printf("[")
				if err := renderTypeDecl(t.Params()[0], w); err != nil {
					return err
				}

				w.Printf("]")
				if err := renderTypeDecl(t.Params()[1], w); err != nil {
					return err
				}

				builtInHandled = true
			case stdlib.List:
				if err := renderTypeDecl(t.Params()[0], w); err != nil {
					return err
				}

				builtInHandled = true
			}

		}

		if !builtInHandled {
			w.Printf("[")
			for i, decl := range t.Params() {
				if err := renderTypeDecl(decl, w); err != nil {
					return err
				}
				if i < len(t.Params())-1 {
					w.Printf(",")
				}
			}
			w.Printf("]")
		}

	case *ast.ChanTypeDeclNode:
		w.Printf("chan ")
		if err := renderTypeDecl(t.TypeDecl(), w); err != nil {
			return err
		}

	case *ast.ArrayTypeDeclNode:
		w.Printf("[%d]", t.SrcArrayTypeDecl().Len())
		if err := renderTypeDecl(t.TypeDecl(), w); err != nil {
			return err
		}

	case *ast.FuncTypeDeclNode:
		w.Printf("func ")
		if err := renderFunc(t.Func(), w); err != nil {
			return err
		}
	default:
		panic("not yet implemented: " + reflect.TypeOf(t).String())
	}

	return nil
}

// visibilityAsName either returns an uppercase name or a lowercase name else.
func visibilityAsName(v src.Visibility, name string) string {
	if len(name) == 0 {
		return name
	}

	switch v {
	case src.Public:
		return MakePublic(name)
	case src.PackagePrivate:
		fallthrough
	case src.Private:
		fallthrough
	case src.Protected:
		return MakePrivate(name)
	default:
		panic("visibility not implemented: " + strconv.Itoa(int(v)))
	}

}
