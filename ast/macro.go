package ast

import (
	"github.com/golangee/src"
	"unicode"
)

// A MacroCtx allows access to actual rendering context but also incorporates the Node less src.MacroCtx.
type MacroCtx interface {
	// Type returns the next parent Type or nil if no such node exists.
	Type() *TypeNode

	// File returns the current SrcFileNode or nil if no such node exists.
	File() *SrcFileNode

	// Node returns the directly enclosing or applicable node, usually a BlockNode.
	Node() Node

	src.MacroCtx
}

// Macro is a function to be called to emit source code directly.
type Macro func(ctx MacroCtx, p func(r ...interface{})) error

// DefaultMacroCtx provides a basic implementation, enough for most use cases.
type DefaultMacroCtx struct {
	node          Node
	resolveStdLib func(name src.Name) src.Name
	importName    func(name src.Name) src.Name
}

// NewDefaultMacroCtx creates a default context.
func NewDefaultMacroCtx(node Node, stdLibResolver func(name src.Name) src.Name, nameImporter func(name src.Name) src.Name) *DefaultMacroCtx {
	return &DefaultMacroCtx{
		node:          node,
		resolveStdLib: stdLibResolver,
		importName:    nameImporter,
	}
}

func (d *DefaultMacroCtx) Type() *TypeNode {
	return ParentTypeNode(d.node)
}

func (d *DefaultMacroCtx) File() *SrcFileNode {
	return ParentSrcFileNode(d.node)
}

func (d *DefaultMacroCtx) Node() Node {
	return d.node
}

func (d *DefaultMacroCtx) MimeType() string {
	f := d.File()
	if f != nil {
		return f.MimeType()
	}

	return ""
}

func (d *DefaultMacroCtx) StdLib(name src.Name) src.Name {
	return d.resolveStdLib(name)
}

func (d *DefaultMacroCtx) Import(name src.Name) src.Name {
	return d.importName(name)
}

func (d *DefaultMacroCtx) Receiver() string {
	fun, isFun := d.Node().(*FuncNode)

	if d.MimeType() == src.Java {
		if isFun && fun.SrcFunc().Static() {
			return ""
		}

		return "this"
	}

	if isFun && fun.SrcFunc().RecName() != "" {
		return fun.SrcFunc().RecName()
	}

	if d.Type() != nil && d.Type().SrcNamedType().Name() != "" {
		// conversion is "correct", because neither Java or Go allow multibyte chars in identifiers
		return string(unicode.ToUpper(rune(d.Type().SrcNamedType().Name()[0])))
	}

	return ""
}
