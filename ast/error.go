package ast

// An Error is a sealed type of a finite set of enumerable and instantiable types.
// Go:
//   create private structs each implementing the error interface and methods named after GroupName and each ErrorCase.
// Java:
//   model as sealed class or (checked) exception?
//
// TODO I don't know what is better: concrete types like this or a macro?
type Error struct {
	GroupName       string // GroupName denotes the actual name of the sealed type set of errors.
	ErrorVisibility Visibility
	Cases           []*ErrorCase // Cases declares possible error cases.
	Obj
}

func NewError(groupName string) *Error {
	return &Error{
		GroupName: groupName,
	}
}

// SetComment sets the nodes comment.
func (n *Error) SetComment(text string) *Error {
	n.ObjComment = NewComment(text)
	n.ObjComment.SetParent(n)
	return n
}

func (n *Error) Visibility() Visibility {
	return n.ErrorVisibility
}

func (n *Error) Identifier() string {
	return n.GroupName
}

func (n *Error) sealedNamedType() {
	panic("implement me")
}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (n *Error) Children() []Node {
	tmp := make([]Node, 0, len(n.Cases))
	for _, enumCase := range n.Cases {
		tmp = append(tmp, enumCase)
	}

	return tmp
}

func (n *Error) AddErrorCases(cases ...*ErrorCase) *Error {
	for _, e := range cases {
		assertNotAttached(e)
		assertSettableParent(e).SetParent(n)
		n.Cases = append(n.Cases, e)
	}

	return n
}

// An ErrorCase declares a unique case of the enumeration.
type ErrorCase struct {
	TypeName        string
	ErrorVisibility Visibility
	Properties      []*Property // Properties are usually reflected as Fields and their according getter-method set.
	Obj
}

func NewErrorCase(name string) *ErrorCase {
	return &ErrorCase{TypeName: name, ErrorVisibility: Public}
}

// SetComment sets the nodes comment.
func (n *ErrorCase) SetComment(text string) *ErrorCase {
	n.ObjComment = NewComment(text)
	n.ObjComment.SetParent(n)
	return n
}

func (n *ErrorCase) ParentError() *Error {
	if e, ok := n.Parent().(*Error); ok {
		return e
	}

	return nil
}

func (n *ErrorCase) Name() string {
	return n.TypeName
}

func (n *ErrorCase) Identifier() string {
	return n.TypeName
}

func (n *ErrorCase) Visibility() Visibility {
	return n.ErrorVisibility
}

func (n *ErrorCase) AddProperties(p ...*Property) *ErrorCase {
	for _, property := range p {
		assertNotAttached(property)
		assertSettableParent(property).SetParent(n)
		n.Properties = append(n.Properties, property)
	}

	return n
}
