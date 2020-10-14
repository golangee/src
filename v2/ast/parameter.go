package ast

import "github.com/golangee/src/v2"

// ParameterNode wraps a src.Param
type ParameterNode struct {
	parent   Node
	srcParam *src.Param
	typeDecl TypeDeclNode
	*payload
}

// NewParameterNode creates a new parameter.
func NewParameterNode(parent Node, srcParam *src.Param) *ParameterNode {
	n := &ParameterNode{
		parent:   parent,
		srcParam: srcParam,
		payload:  newPayload(),
	}

	n.typeDecl = NewTypeDeclNode(n, srcParam.TypeDecl())

	return n
}

// SrcParameter returns the original parameter.
func (p *ParameterNode) SrcParameter() *src.Param {
	return p.srcParam
}

// TypeDecl returns the according type declaration.
func (p *ParameterNode) TypeDecl() TypeDeclNode {
	return p.typeDecl
}

// Parent returns the parent node or nil, if it is the root of the tree.
func (p *ParameterNode) Parent() Node {
	return p.parent
}
