package ast

import (
	"github.com/golangee/src/v2"
	"reflect"
)

// A TypeDeclNode represents a type declaration, not definition, which you find in TypeNode. There is a wrapper
// for each of the 7 src.TypeDecl variants.
type TypeDeclNode interface {
	Node
	sealedTypeDeclNode()
}

// NewTypeDeclNode wraps the given instance and creates a sub tree with parent/children relations to
// create a foundation for context-aware renderers.
func NewTypeDeclNode(parent Node, decl src.TypeDecl) TypeDeclNode {
	switch t := decl.(type) {
	case *src.SimpleTypeDecl:
		return NewSimpleTypeDeclNode(parent, t)
	case *src.TypeDeclPtr:
		return NewTypeDeclPtrNode(parent, t)
	case *src.SliceTypeDecl:
		return NewSliceTypeDeclNode(parent, t)
	case *src.ChanTypeDecl:
		return NewChanTypeDeclNode(parent, t)
	case *src.ArrayTypeDecl:
		return NewArrayTypeDeclNode(parent, t)
	case *src.GenericTypeDecl:
		return NewGenericTypeDeclNode(parent, t)
	case *src.FuncTypeDecl:
		return NewFuncTypeDeclNode(parent, t)
	default:
		panic("not yet implemented: " + reflect.TypeOf(t).String())
	}
}

// ========

// SimpleTypeDeclNode wraps the src.SimpleTypeDecl.
type SimpleTypeDeclNode struct {
	parent            Node
	srcSimpleTypeDecl *src.SimpleTypeDecl
	*payload
}

// NewSimpleTypeDeclNode wraps the src.SimpleTypeDecl.
func NewSimpleTypeDeclNode(parent Node, typeDecl *src.SimpleTypeDecl) *SimpleTypeDeclNode {
	return &SimpleTypeDeclNode{
		parent:            parent,
		srcSimpleTypeDecl: typeDecl,
		payload:           newPayload(),
	}
}

// SrcSimpleTypeDecl returns the original declaration.
func (n *SimpleTypeDeclNode) SrcSimpleTypeDecl() *src.SimpleTypeDecl {
	return n.srcSimpleTypeDecl
}

// Parent returns the parent node or nil, if it is the root of the tree.
func (n *SimpleTypeDeclNode) Parent() Node {
	return n.parent
}

// sealedTypeDeclNode enforces that no others can implement this interface.
func (n *SimpleTypeDeclNode) sealedTypeDeclNode() {
	panic("sealed type")
}

// ========

// TypeDeclPtrNode wraps the src.TypeDeclPtr.
type TypeDeclPtrNode struct {
	parent         Node
	srcTypeDeclPtr *src.TypeDeclPtr
	typeDecl       TypeDeclNode
	*payload
}

// NewSimpleTypeDeclNode wraps the src.SimpleTypeDecl.
func NewTypeDeclPtrNode(parent Node, typeDecl *src.TypeDeclPtr) *TypeDeclPtrNode {
	n := &TypeDeclPtrNode{
		parent:         parent,
		srcTypeDeclPtr: typeDecl,
		payload:        newPayload(),
	}

	n.typeDecl = NewTypeDeclNode(n, typeDecl.TypeDecl())
	return n
}

// TypeDecl returns the wrapped declaration.
func (n *TypeDeclPtrNode) TypeDecl() TypeDeclNode {
	return n.typeDecl
}

// SrcTypeDeclPtr returns the original declaration.
func (n *TypeDeclPtrNode) SrcTypeDeclPtr() *src.TypeDeclPtr {
	return n.srcTypeDeclPtr
}

// Parent returns the parent node or nil, if it is the root of the tree.
func (n *TypeDeclPtrNode) Parent() Node {
	return n.parent
}

// sealedTypeDeclNode enforces that no others can implement this interface.
func (n *TypeDeclPtrNode) sealedTypeDeclNode() {
	panic("sealed type")
}

// ========

// SliceTypeDeclNode wraps the src.SliceTypeDecl.
type SliceTypeDeclNode struct {
	parent           Node
	srcSliceTypeDecl *src.SliceTypeDecl
	typeDecl         TypeDeclNode
	*payload
}

// NewSimpleTypeDeclNode wraps the src.SimpleTypeDecl.
func NewSliceTypeDeclNode(parent Node, typeDecl *src.SliceTypeDecl) *SliceTypeDeclNode {
	n := &SliceTypeDeclNode{
		parent:           parent,
		srcSliceTypeDecl: typeDecl,
		payload:          newPayload(),
	}

	n.typeDecl = NewTypeDeclNode(n, typeDecl.TypeDecl())
	return n
}

// TypeDecl returns the wrapped declaration.
func (n *SliceTypeDeclNode) TypeDecl() TypeDeclNode {
	return n.typeDecl
}

// SrcSliceTypeDecl returns the original declaration.
func (n *SliceTypeDeclNode) SrcSliceTypeDecl() *src.SliceTypeDecl {
	return n.srcSliceTypeDecl
}

// Parent returns the parent node or nil, if it is the root of the tree.
func (n *SliceTypeDeclNode) Parent() Node {
	return n.parent
}

// sealedTypeDeclNode enforces that no others can implement this interface.
func (n *SliceTypeDeclNode) sealedTypeDeclNode() {
	panic("sealed type")
}

// ========

// ArrayTypeDeclNode wraps the src.ArrayTypeDecl.
type ArrayTypeDeclNode struct {
	parent           Node
	srcArrayTypeDecl *src.ArrayTypeDecl
	typeDecl         TypeDeclNode
	*payload
}

// NewArrayTypeDeclNode wraps the src.SimpleTypeDecl.
func NewArrayTypeDeclNode(parent Node, typeDecl *src.ArrayTypeDecl) *ArrayTypeDeclNode {
	n := &ArrayTypeDeclNode{
		parent:           parent,
		srcArrayTypeDecl: typeDecl,
		payload:          newPayload(),
	}

	n.typeDecl = NewTypeDeclNode(n, typeDecl.TypeDecl())
	return n
}

// TypeDecl returns the wrapped declaration.
func (n *ArrayTypeDeclNode) TypeDecl() TypeDeclNode {
	return n.typeDecl
}

// SrcArrayTypeDecl returns the original declaration.
func (n *ArrayTypeDeclNode) SrcArrayTypeDecl() *src.ArrayTypeDecl {
	return n.srcArrayTypeDecl
}

// Parent returns the parent node or nil, if it is the root of the tree.
func (n *ArrayTypeDeclNode) Parent() Node {
	return n.parent
}

// sealedTypeDeclNode enforces that no others can implement this interface.
func (n *ArrayTypeDeclNode) sealedTypeDeclNode() {
	panic("sealed type")
}

// ========

// ChanTypeDeclNode wraps the src.ChanTypeDecl.
type ChanTypeDeclNode struct {
	parent          Node
	srcChanTypeDecl *src.ChanTypeDecl
	typeDecl        TypeDeclNode
	*payload
}

// NewChanTypeDeclNode wraps the src.ChanTypeDecl.
func NewChanTypeDeclNode(parent Node, typeDecl *src.ChanTypeDecl) *ChanTypeDeclNode {
	n := &ChanTypeDeclNode{
		parent:          parent,
		srcChanTypeDecl: typeDecl,
		payload:         newPayload(),
	}

	n.typeDecl = NewTypeDeclNode(n, typeDecl.TypeDecl())
	return n
}

// TypeDecl returns the wrapped declaration.
func (n *ChanTypeDeclNode) TypeDecl() TypeDeclNode {
	return n.typeDecl
}

// SrcChanTypeDecl returns the original declaration.
func (n *ChanTypeDeclNode) SrcChanTypeDecl() *src.ChanTypeDecl {
	return n.srcChanTypeDecl
}

// Parent returns the parent node or nil, if it is the root of the tree.
func (n *ChanTypeDeclNode) Parent() Node {
	return n.parent
}

// sealedTypeDeclNode enforces that no others can implement this interface.
func (n *ChanTypeDeclNode) sealedTypeDeclNode() {
	panic("sealed type")
}

// ========

// GenericTypeDeclNode wraps the src.GenericTypeDecl.
type GenericTypeDeclNode struct {
	parent             Node
	srcGenericTypeDecl *src.GenericTypeDecl
	typeDecl           TypeDeclNode
	params             []TypeDeclNode
	*payload
}

// NewGenericTypeDeclNode wraps the src.GenericTypeDecl.
func NewGenericTypeDeclNode(parent Node, typeDecl *src.GenericTypeDecl) *GenericTypeDeclNode {
	n := &GenericTypeDeclNode{
		parent:             parent,
		srcGenericTypeDecl: typeDecl,
		payload:            newPayload(),
	}

	n.typeDecl = NewTypeDeclNode(n, typeDecl.TypeDecl())
	for _, decl := range typeDecl.Params() {
		n.params = append(n.params, NewTypeDeclNode(n, decl))
	}

	return n
}

// TypeDecl returns the wrapped base type declaration which is parameterized by Params.
func (n *GenericTypeDeclNode) TypeDecl() TypeDeclNode {
	return n.typeDecl
}

// TypeDecl returns the wrapped type parameter declarations.
func (n *GenericTypeDeclNode) Params() []TypeDeclNode {
	return n.params
}

// SrcArrayTypeDecl returns the original declaration.
func (n *GenericTypeDeclNode) SrcGenericTypeDecl() *src.GenericTypeDecl {
	return n.srcGenericTypeDecl
}

// Parent returns the parent node or nil, if it is the root of the tree.
func (n *GenericTypeDeclNode) Parent() Node {
	return n.parent
}

// sealedTypeDeclNode enforces that no others can implement this interface.
func (n *GenericTypeDeclNode) sealedTypeDeclNode() {
	panic("sealed type")
}

// ========

// FuncTypeDeclNode wraps the src.FuncTypeDecl.
type FuncTypeDeclNode struct {
	parent          Node
	srcFuncTypeDecl *src.FuncTypeDecl
	typeDecl        TypeDeclNode
	params          []*ParameterNode
	results         []*ParameterNode
	*payload
}

// NewGenericTypeDeclNode wraps the src.FuncTypeDecl.
func NewFuncTypeDeclNode(parent Node, typeDecl *src.FuncTypeDecl) *FuncTypeDeclNode {
	n := &FuncTypeDeclNode{
		parent:          parent,
		srcFuncTypeDecl: typeDecl,
		payload:         newPayload(),
	}

	for _, decl := range typeDecl.InputParams() {
		n.params = append(n.params, NewParameterNode(n, decl))
	}

	for _, decl := range typeDecl.OutputParams() {
		n.params = append(n.params, NewParameterNode(n, decl))
	}

	return n
}

// SrcFuncTypeDecl returns the original declaration.
func (n *FuncTypeDeclNode) SrcFuncTypeDecl() *src.FuncTypeDecl {
	return n.srcFuncTypeDecl
}

// TypeDecl returns the wrapped base type declaration which is parameterized by Params.
func (n *FuncTypeDeclNode) TypeDecl() TypeDeclNode {
	return n.typeDecl
}

// InputParams returns the wrapped input parameter declarations.
func (n *FuncTypeDeclNode) InputParams() []*ParameterNode {
	return n.params
}

// OutputParams returns the wrapped input parameter declarations.
func (n *FuncTypeDeclNode) OutputParams() []*ParameterNode {
	return n.results
}

// Parent returns the parent node or nil, if it is the root of the tree.
func (n *FuncTypeDeclNode) Parent() Node {
	return n.parent
}

// sealedTypeDeclNode enforces that no others can implement this interface.
func (n *FuncTypeDeclNode) sealedTypeDeclNode() {
	panic("sealed type")
}
