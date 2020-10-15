package ast

import "github.com/golangee/src/v2"

// A StructNode represents a data class (e.g. PoJo) or record. However, the actual expected semantic is
// the Go semantic, other languages which have no value/reference expression, will have a problem and need
// to fallback to their nearest idiomatic representation.
type StructNode struct {
	parent      *TypeNode
	srcStruct   *src.Struct
	fields      []*FieldNode
	annotations []*AnnotationNode
	types       []*TypeNode
	methods     []*FuncNode
	*payload
}

// NewStructNode wraps the given instance and creates a sub tree with parent/children relations to
// create a foundation for context-aware renderers.
func NewStructNode(parent *TypeNode, srcStruct *src.Struct) *StructNode {
	n := &StructNode{
		parent:    parent,
		srcStruct: srcStruct,
		payload:   newPayload(),
	}

	for _, f := range srcStruct.Methods() {
		n.methods = append(n.methods, NewFuncNode(n, f))
	}

	for _, field := range srcStruct.Fields() {
		n.fields = append(n.fields, NewFieldNode(n, field))
	}

	for _, annotation := range srcStruct.Annotations() {
		n.annotations = append(n.annotations, NewAnnotationNode(n, annotation))
	}

	for _, namedType := range srcStruct.Types() {
		n.types = append(n.types, NewTypeNode(n, namedType))
	}

	return n
}

// Name returns the declared identifier which must be unique per package.
func (n *StructNode) Name() string {
	return n.srcStruct.Name()
}

// Doc returns the package documentation.
func (n *StructNode) Doc() string {
	return n.srcStruct.Doc()
}

// SrcStruct returns the original struct.
func (n *StructNode) SrcStruct() *src.Struct {
	return n.srcStruct
}

// Fields returns the backing slice of the wrapped fields.
func (n *StructNode) Fields() []*FieldNode {
	return n.fields
}

// Parent returns the parent node or nil, if it is the root of the tree.
func (n *StructNode) Parent() Node {
	return n.parent
}

// Annotations returns all registered annotations.
func (n *StructNode) Annotations() []*AnnotationNode {
	return n.annotations
}

// Types returns all defines subtypes in the scope of this struct.
func (n *StructNode) Types() []*TypeNode {
	return n.types
}

// Methods returns the backing slice of the wrapped methods.
func (n *StructNode) Methods() []*FuncNode {
	return n.methods
}
