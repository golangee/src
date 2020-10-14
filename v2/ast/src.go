package ast

import "github.com/golangee/src/v2"

// SrcFileNode represents a compilation unit.
type SrcFileNode struct {
	parent  *PkgNode
	srcFile *src.SrcFile
	types   []*TypeNode
	*payload
}

// NewSrcFileNode wraps the given instance and creates a sub tree with parent/children relations to
// create a foundation for context-aware renderers.
func NewSrcFileNode(parent *PkgNode, file *src.SrcFile) *SrcFileNode {
	n := &SrcFileNode{
		parent:  parent,
		srcFile: file,
		payload: newPayload(),
	}

	for _, namedType := range file.Types() {
		n.types = append(n.types, NewTypeNode(n, namedType))
	}

	return n
}

// PkgNode walks up the hierarchy and returns the containing package.
func (n *SrcFileNode) PkgNode() *PkgNode {
	return n.parent
}

// SrcFile returns the wrapped type.
func (n *SrcFileNode) SrcFile() *src.SrcFile {
	return n.srcFile
}

// Types returns the wrapped node types.
func (n *SrcFileNode) Types() []*TypeNode {
	return n.types
}

// Parent returns the parent node or nil, if it is the root of the tree.
func (n *SrcFileNode) Parent() Node {
	return n.parent
}
