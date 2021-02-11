package ast

// An Enum is an enumerable type. Its natural form is an integer but it may also include a string representation.
type Enum struct {
	TypeName   string // TypeName denotes the actual name of this type.
	BaseType   Name   // BaseType must be a primitive like int or string in Go.
	Implements []Name // Implements denotes a bunch of interfaces which must be implemented by this struct. Depending on the renderer (like Go) this has no effect.
	Cases      []*EnumCase
	Obj
}

func (n *Enum) Identifier() string {
	return n.TypeName
}

func (n *Enum) sealedNamedType() {
	panic("implement me")
}

// An EnumCase declares a unique case of the enumeration.
type EnumCase struct {
	TypeName string
	Value    *BasicLit
	Obj
}

func (n *EnumCase) Name() string {
	return n.TypeName
}
