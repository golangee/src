package ast2

import "strconv"

// A Pos describes a resolved position within a file.
type Pos struct {
	// File contains the absolute file path.
	File string
	// Line denotes the one-based line number In the denoted File.
	Line int
	// Col denotes the one-based column number In the denoted Line.
	Col int
}

// String returns the content In the "file:line:col" format.
func (p Pos) String() string {
	return p.File + ":" + strconv.Itoa(p.Line) + ":" + strconv.Itoa(p.Col)
}

// Tags is just a simple string/interface map to store arbitrary Values. This is especially useful
// to attach hidden generator information, which otherwise do not fit into an AST.
type Tags map[string]interface{}

// Get returns the according value or nil. This is nil safe.
func (t Tags) Get(key string) interface{} {
	if t == nil {
		return nil
	}

	return t[key]
}

// A Node represents the common contract
type Node interface {
	// Pos returns the actual starting position of this Node.
	Pos() Pos

	// End is the position of the first char after the node.
	End() Pos

	// Parent returns the parent Node or nil if undefined. This recursive implementation may be considered as
	// unnecessary and even as an anti pattern within an AST but the core feature is to perform semantic validations
	// which requires a lot of down/up iterations through the (entire) AST. Keeping the relational relation
	// at the node level keeps things simple and we don't need to pass (path) contexts everywhere.
	Parent() Node

	// Tags returns access to arbitrary tags. May be nil, so always use the Get accessor.
	Tags() Tags

	// Comment returns an optional comment node.
	Comment() *Comment
}

// A Parent is a Node and may contain other nodes as children. This is used to simplify algorithms based on Walk.
type Parent interface {
	Node
	// Children returns a defensive copy of the underlying slice. However the Node references are shared.
	Children() []Node
}

// Obj is actually a helper to implement a Node by embedding the Obj
type Obj struct {
	ObjPos     Pos
	ObjEnd     Pos
	ObjParent  Node
	ObjTags    Tags
	ObjComment *Comment // the actual comment of the logical object
}

func (n *Obj) Pos() Pos {
	return n.ObjPos
}

func (n *Obj) End() Pos {
	return n.ObjEnd
}

func (n *Obj) Parent() Node {
	return n.ObjParent
}

func (n *Obj) SetParent(p Node) {
	n.ObjParent = p
}

func (n *Obj) Tags() Tags {
	return n.ObjTags
}

func (n *Obj) Comment() *Comment {
	return n.ObjComment
}

type SettableParent interface {
	Node
	SetParent(p Node)
}
