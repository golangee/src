package ast

// Node represents a declared element in the abstract syntax tree.
type Node interface {
	// Parent returns the parent node, or nil if its the root node.
	Parent() Node

	// Value returns nil or the associated payload value.
	Value(key interface{}) interface{}

	// SetValue updates the payload value for the given key.
	SetValue(key, value interface{})
}
