package ast

import (
	"github.com/golangee/src"
)

// An InterfaceNode represents a (named) interface type.
type InterfaceNode struct {
	parent       Node
	srcInterface *src.Interface
	methods      []*FuncNode
	annotations  []*AnnotationNode
	types        []*TypeNode
	*payload
}

// NewInterfaceNode wraps the given instance and creates a sub tree with parent/children relations to
// create a foundation for context-aware renderers.
func NewInterfaceNode(parent Node, iface *src.Interface) *InterfaceNode {
	n := &InterfaceNode{
		parent:       parent,
		srcInterface: iface,
		payload:      newPayload(),
	}

	for _, f := range iface.Methods() {
		n.methods = append(n.methods, NewFuncNode(n, f))
	}

	for _, annotation := range iface.Annotations() {
		n.annotations = append(n.annotations, NewAnnotationNode(n, annotation))
	}

	for _, namedType := range iface.Types() {
		n.types = append(n.types, NewTypeNode(n, namedType))
	}

	return n
}

// Name returns the declared identifier which must be unique per package.
func (n *InterfaceNode) Name() string {
	return n.srcInterface.Name()
}

// Doc returns the package documentation.
func (n *InterfaceNode) Doc() string {
	return n.srcInterface.Doc()
}

// SrcInterface returns the original interface.
func (n *InterfaceNode) SrcInterface() *src.Interface {
	return n.srcInterface
}

// Methods returns the backing slice of the wrapped methods.
func (n *InterfaceNode) Methods() []*FuncNode {
	return n.methods
}

// Parent returns the parent node or nil, if it is the root of the tree.
func (n *InterfaceNode) Parent() Node {
	return n.parent
}

// Annotations returns all registered annotations.
func (n *InterfaceNode) Annotations() []*AnnotationNode {
	return n.annotations
}

// Types returns all defines subtypes in the scope of this interface.
func (n *InterfaceNode) Types() []*TypeNode {
	return n.types
}
