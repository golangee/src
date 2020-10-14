package ast

import "github.com/golangee/src/v2"

// FuncNode represents a method or function, depending on the context.
type FuncNode struct {
	parent  Node
	srcFunc *src.Func
	*payload
}

// NewFuncNode wraps the given instance and creates a sub tree with parent/children relations to
// create a foundation for context-aware renderers.
func NewFuncNode(parent Node, fun *src.Func) *FuncNode {
	return &FuncNode{
		parent:  parent,
		srcFunc: fun,
		payload: newPayload(),
	}
}

// SrcFunc returns the original func.
func (n *FuncNode) SrcFunc() *src.Func {
	return n.srcFunc
}

// Parent returns the parent node or nil, if it is the root of the tree.
func (n *FuncNode) Parent() Node {
	return n.parent
}
