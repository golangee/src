package ast

import "github.com/golangee/src/v2"

// An InterfaceNode represents a (named) interface type.
type InterfaceNode struct {
	parent       Node
	srcInterface *src.Interface
	methods      []*FuncNode
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

	return n
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
