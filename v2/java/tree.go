package java

import (
	"github.com/golangee/src/v2"
	"reflect"
)

type node interface {
	// Parent returns the parent node, or nil if its the root node.
	Parent() node
}

// =====

type modNode struct {
	srcModule *src.Module
	packages  []*pkgNode
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

// Parent returns the parent node or nil, if it is the root of the tree.
func (n *modNode) Parent() node {
	return nil
}

// =====

type pkgNode struct {
	parent     *modNode
	srcPackage *src.Package
	files      []*srcFileNode
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

// Parent returns the parent node or nil, if it is the root of the tree.
func (n *pkgNode) Parent() node {
	return n.parent
}

// =====

type srcFileNode struct {
	parent   *pkgNode
	srcFile  *src.SrcFile
	types    []*typeNode
	importer *importer
}

func newSrcFileNode(parent *pkgNode, file *src.SrcFile) *srcFileNode {
	n := &srcFileNode{
		parent:   parent,
		srcFile:  file,
		importer: newImporter(),
	}

	for _, namedType := range file.Types() {
		n.types = append(n.types, newTypeNode(n, namedType))
	}

	return n
}

// Parent returns the parent node or nil, if it is the root of the tree.
func (n *srcFileNode) Parent() node {
	return n.parent
}

// importName returns the (optional) shorter import name. An internal state is
// created to ensure, that
func (n *srcFileNode) importName(name src.Name) string {
	panic("todo")
}

// =====

type typeNode struct {
	parent       *srcFileNode
	srcNamedType src.NamedType
	namedNode    interface{} // one of *structNode
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

// Parent returns the parent node or nil, if it is the root of the tree.
func (n *typeNode) Parent() node {
	return n.parent
}

// =====

type structNode struct {
	parent    *typeNode
	srcStruct *src.Struct
	fields    []*fieldNode
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

// Parent returns the parent node or nil, if it is the root of the tree.
func (n *structNode) Parent() node {
	return n.parent
}

// =====

type fieldNode struct {
	parent       *structNode
	srcField     *src.Field
	typeDeclNode *typeDeclNode
	annotations  []*annotationNode
}

func newFieldNode(parent *structNode, field *src.Field) *fieldNode {
	n := &fieldNode{
		parent:   parent,
		srcField: field,
	}

	n.typeDeclNode = newTypeDeclNode(n, field.TypeDecl())
	for _, annotation := range field.Annotations() {
		n.annotations = append(n.annotations, newAnnotationNode(n, annotation))
	}

	return n
}

// Parent returns the parent node or nil, if it is the root of the tree.
func (n *fieldNode) Parent() node {
	return n.parent
}

// =====

type typeDeclNode struct {
	parent      node
	srcTypeDecl src.TypeDecl
}

func newTypeDeclNode(parent node, decl src.TypeDecl) *typeDeclNode {
	return &typeDeclNode{
		parent:      parent,
		srcTypeDecl: decl,
	}
}

// Parent returns the parent node or nil, if it is the root of the tree.
func (n *typeDeclNode) Parent() node {
	return n.parent
}

// =====

type annotationNode struct {
	parent        node
	srcAnnotation *src.Annotation
}

func newAnnotationNode(parent node, a *src.Annotation) *annotationNode {
	return &annotationNode{
		parent:        parent,
		srcAnnotation: a,
	}
}

// Parent returns the parent node or nil, if it is the root of the tree.
func (n *annotationNode) Parent() node {
	return n.parent
}
