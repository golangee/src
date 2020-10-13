package java

import (
	"github.com/golangee/src/v2"
	"reflect"
)

type modNode struct {
	srcModule *src.Module
	packages  []*pkgNode
}

type pkgNode struct {
	parent     *modNode
	srcPackage *src.Package
	files      []*srcFileNode
}

type importNode struct {
	parent     *srcFileNode
	name       src.Name // something like my.package.name.MyOuterClass.MyInnerClass
	qualifier  string   // something like my.package.name
	identifier src.Name // may be still something like MyOuterClass.MyInnerClass
	wildcard   bool     // indicates a wildcard import
}

type srcFileNode struct {
	parent  *pkgNode
	srcFile *src.SrcFile
	types   []*typeNode
	imports []importNode
}

func newSrcFileNode(parent *pkgNode, file *src.SrcFile) *srcFileNode {
	n := &srcFileNode{
		parent:  parent,
		srcFile: file,
	}

	for _, namedType := range file.Types() {
		n.types = append(n.types, newTypeNode(n, namedType))
	}

	return n
}

// importName returns the (optional) shorter import name. An internal state is
// created to ensure, that
func (n *srcFileNode) importName(name src.Name) string {
	panic("todo")
}

type typeNode struct {
	parent       *srcFileNode
	srcNamedType src.NamedType
	namedNode    interface{} // one of *structNode
}

type structNode struct {
	parent    *typeNode
	srcStruct *src.Struct
	fields    []*fieldNode
}

type fieldNode struct {
	parent       *structNode
	srcField     *src.Field
	typeDeclNode *typeDeclNode
}

type typeDeclNode struct {
	parent      interface{}
	srcTypeDecl src.TypeDecl
}

func newModNode(mod *src.Module) *modNode {
	n := &modNode{
		srcModule: mod,
	}

	for _, p := range mod.Packages() {
		n.packages = append(n.packages, newPkgNode(n, p))
	}

	return n
}

func newPkgNode(parent *modNode, pkg *src.Package) *pkgNode {
	n := &pkgNode{
		parent:     parent,
		srcPackage: pkg,
	}

	for _, file := range pkg.SrcFiles() {
		n.files = append(n.files, newSrcFileNode(n, file))
	}

	return n
}

func newTypeNode(parent *srcFileNode, t src.NamedType) *typeNode {
	n := &typeNode{
		parent:       parent,
		srcNamedType: t,
	}

	switch t := t.(type) {
	case *src.Struct:
		n.namedNode = newStructNode(n, t)
	default:
		panic("not yet implemented: " + reflect.TypeOf(t).String())
	}

	return n
}

func newStructNode(parent *typeNode, srcStruct *src.Struct) *structNode {
	n := &structNode{
		parent:    parent,
		srcStruct: srcStruct,
	}

	for _, field := range srcStruct.Fields() {
		n.fields = append(n.fields, newFieldNode(n, field))
	}

	return n
}

func newFieldNode(parent *structNode, field *src.Field) *fieldNode {
	n := &fieldNode{
		parent:   parent,
		srcField: field,
	}

	n.typeDeclNode = newTypeDeclNode(n, field.TypeDecl())
	return n
}

func newTypeDeclNode(parent interface{}, decl src.TypeDecl) *typeDeclNode {
	return &typeDeclNode{
		parent:      parent,
		srcTypeDecl: decl,
	}
}
