package ast

import (
	"github.com/golangee/src/v2"
	"reflect"
)

// A TypeNode represents a (named) type declaration. Nameless or anonymous types may be legal.
type TypeNode struct {
	parent       *SrcFileNode
	srcNamedType src.NamedType
	namedNode    Node
	*payload
}

// NewTypeNode wraps the given instance and creates a sub tree with parent/children relations to
// create a foundation for context-aware renderers.
func NewTypeNode(parent *SrcFileNode, t src.NamedType) *TypeNode {
	n := &TypeNode{
		parent:       parent,
		srcNamedType: t,
		payload:      newPayload(),
	}

	switch t := t.(type) {
	case *src.Struct:
		n.namedNode = NewStructNode(n, t)
	case *src.Interface:
		n.namedNode = NewInterfaceNode(n, t)
	default:
		panic("not yet implemented: " + reflect.TypeOf(t).String())
	}

	return n
}

// SrcNamedType returns the original named type.
func (n *TypeNode) SrcNamedType() src.NamedType {
	return n.srcNamedType
}

// NamedNode returns the concrete wrapped node. This is one of
//  *StructNode
//  *InterfaceNode
func (n *TypeNode) NamedNode() Node {
	return n.namedNode
}

// Parent returns the parent node or nil, if it is the root of the tree.
func (n *TypeNode) Parent() Node {
	return n.parent
}
