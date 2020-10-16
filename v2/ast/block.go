package ast

import "github.com/golangee/src/v2"

// BlockNode represents a Block, which is mostly a raw piece of text, scattered with other Blocks, Names and TypeDecls.
type BlockNode struct {
	parent   *FuncNode
	srcBlock *src.Block
	*payload
}

// NewBlockNode wraps the given instance and creates a sub tree with parent/children relations to
// create a foundation for context-aware renderers.
func NewBlockNode(parent *FuncNode, block *src.Block) *BlockNode {
	n := &BlockNode{
		parent:   parent,
		srcBlock: block,
		payload:  newPayload(),
	}

	return n
}

// SrcBlock returns the original block.
func (n *BlockNode) SrcBlock() *src.Block {
	return n.srcBlock
}

// Parent returns the parent node or nil, if it is the root of the tree.
func (n *BlockNode) Parent() Node {
	return n.parent
}

// Eval will recursively execute and transform all Macro and src.Macro and src.TypeDecl.
func (n *BlockNode) Eval(stdLibResolver func(name src.Name) src.Name, nameImporter func(name src.Name) src.Name) ([]interface{}, error) {
	// TODO
	return nil, nil
}
