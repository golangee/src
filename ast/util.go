package ast

// ParentSrcFileNode searches from leaf to root and returns nil if no SrcFileNode has been found.
func ParentSrcFileNode(node Node) *SrcFileNode {
	if node == nil {
		return nil
	}

	if fnode, ok := node.(*SrcFileNode); ok {
		return fnode
	}

	return ParentSrcFileNode(node.Parent())
}

// ParentTypeNode searches from leaf to root and returns nil if no SrcFileNode has been found.
func ParentTypeNode(node Node) *TypeNode {
	if node == nil {
		return nil
	}

	if fnode, ok := node.(*TypeNode); ok {
		return fnode
	}

	return ParentTypeNode(node.Parent())
}
