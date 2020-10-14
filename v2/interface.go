package src

var _ NamedType = (*Interface)(nil)

// An Interface is a contract which defines a method set and allows polymorphism without inheritance. If a
// body is declared, it depends on the actual renderer, if a default method will be emitted (e.g. for Java).
// Go does not support default methods in interfaces.
type Interface struct {
	doc        string
	name       string
	visibility Visibility
	methods    []*Func
}

// NewStruct returns a new named struct type. A struct is always mutable, but may be used either in a value
// or pointer context. Structs are straightforward in Go but in Java just a PoJo. We do not use records, because
// they have a different semantic (read only).
func NewInterface(name string) *Interface {
	return &Interface{name: name}
}

// Name returns the declared identifier which must be unique per package.
func (s *Interface) Name() string {
	return s.name
}

// SetName updates the interfaces identifier which must be unique per package.
func (s *Interface) SetName(name string) *Interface {
	s.name = name
	return s
}

func (s *Interface) sealedNamedType() {
	panic("implement me")
}

// SetVisibility sets the visibility. The default is Public.
func (s *Interface) SetVisibility(v Visibility) *Interface {
	s.visibility = v
	return s
}

// Visibility returns the current visibility. The default is Public.
func (s *Interface) Visibility() Visibility {
	return s.visibility
}

// SetDoc sets the package documentation, which is e.g. emitted to a doc.go or a package-info.java.
func (s *Interface) SetDoc(doc string) *Interface {
	s.doc = doc
	return s
}

// Doc returns the package documentation.
func (s *Interface) Doc() string {
	return s.doc
}

// Methods returns all available functions.
func (s *Interface) Methods() []*Func {
	return s.methods
}

// AddMethods appends more methods to this interfaces contract.
func (s *Interface) AddMethods(f ...*Func) *Interface {
	s.methods = append(s.methods, f...)
	return s
}
