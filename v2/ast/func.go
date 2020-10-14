package ast

import "github.com/golangee/src/v2"

// FuncNode represents a method or function, depending on the context.
type FuncNode struct {
	parent  Node
	srcFunc *src.Func
	*payload
	params      []*ParameterNode
	results     []*ParameterNode
	annotations []*AnnotationNode
}

// NewFuncNode wraps the given instance and creates a sub tree with parent/children relations to
// create a foundation for context-aware renderers.
func NewFuncNode(parent Node, fun *src.Func) *FuncNode {
	n := &FuncNode{
		parent:  parent,
		srcFunc: fun,
		payload: newPayload(),
	}

	for _, param := range fun.Params() {
		n.params = append(n.params, NewParameterNode(n, param))
	}

	for _, param := range fun.Results() {
		n.results = append(n.results, NewParameterNode(n, param))
	}

	for _, annotation := range fun.Annotations() {
		n.annotations = append(n.annotations, NewAnnotationNode(n, annotation))
	}
	return n
}

// InputParams returns the wrapped input parameter declarations.
func (n *FuncNode) InputParams() []*ParameterNode {
	return n.params
}

// OutputParams returns the wrapped input parameter declarations.
func (n *FuncNode) OutputParams() []*ParameterNode {
	return n.results
}

// SrcFunc returns the original func.
func (n *FuncNode) SrcFunc() *src.Func {
	return n.srcFunc
}

// Parent returns the parent node or nil, if it is the root of the tree.
func (n *FuncNode) Parent() Node {
	return n.parent
}

// Annotations returns all registered annotations.
func (n *FuncNode) Annotations() []*AnnotationNode {
	return n.annotations
}
