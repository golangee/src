package ast

import (
	"github.com/golangee/src"
)

// A PkgNode represents a package and contains compilation units (source code files).
type PkgNode struct {
	parent     *ModNode
	srcPackage *src.Package
	files      []*SrcFileNode
	*payload
}

// NewPkgNode wraps the given instance and creates a sub tree with parent/children relations to
// create a foundation for context-aware renderers.
func NewPkgNode(parent *ModNode, pkg *src.Package) *PkgNode {
	n := &PkgNode{
		parent:     parent,
		srcPackage: pkg,
		payload:    newPayload(),
	}

	for _, file := range pkg.SrcFiles() {
		n.files = append(n.files, NewSrcFileNode(n, file))
	}

	return n
}

// Files returns the backing slice of the file nodes.
func (n *PkgNode) Files() []*SrcFileNode {
	return n.files
}

// SrcPackage returns the wrapped instance.
func (n *PkgNode) SrcPackage() *src.Package {
	return n.srcPackage
}

// Parent returns the parent node or nil, if it is the root of the tree.
func (n *PkgNode) Parent() Node {
	return n.parent
}
