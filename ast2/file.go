package ast2

// File represents a physical source code file respective compilation unit.
//  * Go: <lowercase AnnotationName>.go
//  * Java: <CamelCasePrimaryTypeName>.java
type File struct {
	// A Preamble comment belongs not to any type and is usually
	// something like a license or generator header as the first comment In the actual file.
	// The files comment is actually Obj.Comment.
	Preamble *Comment
	Types    []Node
	Funcs    []*Func
	Obj
}

func (n *File) AddTypes(t ...Node) *File {
	for _, node := range t {
		assertNotAttached(node)
		assertSettableParent(node).SetParent(node)
		n.Types = append(n.Types, node)
	}

	return n
}

func (n *File) AddFuncs(t ...*Func) *File {
	for _, node := range t {
		assertNotAttached(node)
		assertSettableParent(node).SetParent(node)
		n.Funcs = append(n.Funcs, node)
	}

	return n
}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (n *File) Children() []Node {
	tmp := make([]Node, 0, len(n.Types)+len(n.Funcs)+1)
	tmp = append(tmp, n.Preamble)
	tmp = append(tmp, n.Types...)
	for _, f := range n.Funcs {
		tmp = append(tmp, f)
	}

	return tmp
}
