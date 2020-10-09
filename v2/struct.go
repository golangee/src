package src

var _ NamedType = (*Struct)(nil)

// A Struct is a record type
type Struct struct {
	doc        string
	name       string
	visibility Visibility
}

// NewStruct returns a new named struct type. A struct is always mutable, but may be used either in a value
// or pointer context. Structs are straightforward in Go but in Java just a PoJo. We do not use records, because
// they have a different semantic (read only).
func NewStruct(name string) *Struct {
	return &Struct{name: name}
}

// Name returns the declared identifier which must be unique per package.
func (s *Struct) Name() string {
	return s.name
}

func (s *Struct) isNamedType() {
	panic("implement me")
}

// SetVisibility sets the visibility. The default is Public.
func (s *Struct) SetVisibility(v Visibility) *Struct {
	s.visibility = v
	return s
}

// Visibility returns the current visibility. The default is Public.
func (s *Struct) Visibility() Visibility {
	return s.visibility
}

// SetDoc sets the package documentation, which is e.g. emitted to a doc.go or a package-info.java.
func (s *Struct) SetDoc(doc string) *Struct {
	s.doc = doc
	return s
}

// Doc returns the package documentation.
func (s *Struct) Doc() string {
	return s.doc
}
