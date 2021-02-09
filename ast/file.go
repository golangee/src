package ast

import "reflect"

// File represents a physical source code file respective compilation unit.
//  * Go: <lowercase AnnotationName>.go
//  * Java: <CamelCasePrimaryTypeName>.java
type File struct {
	// A Preamble comment belongs not to any type and is usually
	// something like a license or generator header as the first comment In the actual file.
	// The files comment is actually Obj.Comment.
	Preamble *Comment
	Name     string
	Types    []Node
	Functions    []*Func
	Obj
}

// NewFile allocates a new File.
func NewFile(name string) *File {
	return &File{Name: name}
}

func (n *File) AddTypes(t ...Node) *File {
	for _, node := range t {
		assertNotAttached(node)
		assertSettableParent(node).SetParent(node)
		n.Types = append(n.Types, node)
	}

	return n
}

// Pkg asserts that the parent is a Pkg instance and returns it.
func (n *File) Pkg() *Pkg {
	if p, ok := n.Parent().(*Pkg); ok {
		return p
	}

	panic("expected parent to be a *Pkg, but was: " + reflect.TypeOf(n.Parent()).Name())
}

func (n *File) AddFuncs(t ...*Func) *File {
	for _, node := range t {
		assertNotAttached(node)
		assertSettableParent(node).SetParent(node)
		n.Functions = append(n.Functions, node)
	}

	return n
}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (n *File) Children() []Node {
	tmp := make([]Node, 0, len(n.Types)+len(n.Functions)+1)
	tmp = append(tmp, n.Preamble)
	tmp = append(tmp, n.Types...)
	for _, f := range n.Functions {
		tmp = append(tmp, f)
	}

	return tmp
}
