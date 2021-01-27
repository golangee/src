package ast2

// A Block represents a lexical group of declarations and statements. Usually a block also introduces a scope and
// can also be nested.
//  Go/Java: { ... }
type Block struct {
	Nodes []Node
	Obj
}

func (n *Block) Add(nodes ...Node) *Block {
	for _, node := range nodes {
		n.Nodes = append(n.Nodes, node)
		assertNotAttached(node)
		assertSettableParent(node).SetParent(n)
	}

	return n
}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (n *Block) Children() []Node {
	return append(make([]Node, 0, len(n.Nodes)), n.Nodes...)
}
