package ast

import "github.com/golangee/src/v2"

// A FieldNode represents a struct field.
type FieldNode struct {
	parent       *StructNode
	srcField     *src.Field
	typeDeclNode TypeDeclNode
	annotations  []*AnnotationNode
	*payload
}

// NewFieldNode wraps the given instance and creates a sub tree with parent/children relations to
// create a foundation for context-aware renderers.
func NewFieldNode(parent *StructNode, field *src.Field) *FieldNode {
	n := &FieldNode{
		parent:   parent,
		srcField: field,
		payload:  newPayload(),
	}

	n.typeDeclNode = NewTypeDeclNode(n, field.TypeDecl())
	for _, annotation := range field.Annotations() {
		n.annotations = append(n.annotations, NewAnnotationNode(n, annotation))
	}

	return n
}

// SrcField returns the original field.
func (n *FieldNode) SrcField() *src.Field {
	return n.srcField
}

// TypeDecl returns the wrapped type declaration of the field.
func (n *FieldNode) TypeDecl() TypeDeclNode {
	return n.typeDeclNode
}

// Annotations returns the backing slice of all field annotations.
func (n *FieldNode) Annotations() []*AnnotationNode {
	return n.annotations
}

// Parent returns the parent node or nil, if it is the root of the tree.
func (n *FieldNode) Parent() Node {
	return n.parent
}
