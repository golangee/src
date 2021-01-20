package ast

import (
	"github.com/golangee/src"
)

// A ModNode is the root and describes a module with packages.
type ModNode struct {
	srcModule *src.Module
	packages  []*PkgNode
	*payload
}

// NewModNode wraps the given instance and creates a sub tree with parent/children relations to
// create a foundation for context-aware renderers.
func NewModNode(mod *src.Module) *ModNode {
	n := &ModNode{
		srcModule: mod,
		payload:   newPayload(),
	}

	for _, p := range mod.Packages() {
		n.packages = append(n.packages, NewPkgNode(n, p))
	}

	return n
}

// Packages returns the backing slice of the package nodes.
func (n *ModNode) Packages() []*PkgNode {
	return n.packages
}

// SrcModule returns the wrapped instance.
func (n *ModNode) SrcModule() *src.Module {
	return n.srcModule
}

// Parent returns the parent node or nil, if it is the root of the tree.
func (n *ModNode) Parent() Node {
	return nil
}
