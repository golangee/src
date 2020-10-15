package ast

import "github.com/golangee/src/v2"

// A MacroCtx allows access to actual rendering context.
type MacroCtx interface {
	// Type returns the next parent Type or nil if no such node exists.
	Type() *TypeNode

	// File returns the current SrcFileNode or nil if no such node exists.
	File() *SrcFileNode

	// MimeType is a shortcut to the files mimetype. If not applicable, returns the empty string.
	MimeType() string

	// Node returns the directly enclosing or applicable node, usually a BlockNode.
	Node() Node

	// StdLib tries to resolve the given name into a standard library type defined by the actual render context
	// without importing it.
	StdLib(name src.Name) src.Name

	// Import tries to resolve the given name from the standard library and imports it for the usage in the current
	// file.
	Import(name src.Name) src.Name
}

// Macro is a function to be called to emit source code directly.
type Macro func(ctx MacroCtx, w interface {
	Printf(str string, args ...interface{})
}) error
