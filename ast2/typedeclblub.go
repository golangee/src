package ast2


//TODO is this a TypeDef?
// A TypeDecl binds an Identifier (the type name) to a Type. Actual types are the following Node type:
//  * Go:
//    * AliasDecl
//    * StructDecl
//    * InterfaceDecl
//    * FunctionDecl
//    * TypeDef
//    * EnumDecl (artificial type based on string)
//
//  * Java:
//    * ClassDecl
//    * InterfaceDecl
//    * EnumDecl
//    * AnnotationDecl
type TypeDecl2 struct {
	Visibility Visibility   // things like Public, Protected or Private
	Identifier string       // Identifier must be unique in the current scope, usually at the package or block level
	Params     []*TypeParam // the introduced type parameter
	Type       Node         // the actual type, see above
	Obj
}

// SetType attaches the given node as the actual type.
func (n *TypeDecl2) SetType(t Node) *TypeDecl2 {
	assertNotAttached(t)
	assertSettableParent(t).SetParent(n)

	return n
}

// Children always returns the Type. A nil Type will raise a panic.
func (n *TypeDecl2) Children() []Node {
	if n.Type == nil {
		panic("invalid type declaration")
	}

	return []Node{n.Type}
}
