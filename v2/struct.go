package src

var _ NamedType = (*Struct)(nil)

// A Struct is actually a data type, like a record. Depending on the language, it can be used in a value or reference
// context. If supported, the primary use case should be the usage as a value to improve conclusiveness and
// performance by avoiding heap allocation (and potentially GC overhead). Inheritance is not possible, but other
// types may be embedded (e.g. in Go). Languages like Java use just simple classes (PoJos), because records have no
// real use (they are just syntax sugar for a class with final members). In contrast to that, Go cannot express
// final fields.
type Struct struct {
	doc        string
	name       string
	visibility Visibility
	fields     []*Field
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

// AddFields appends the given fields to the struct.
func (s *Struct) AddFields(fields ...*Field) *Struct {
	s.fields = fields
	return s
}

// Fields returns the currently configured fields.
func (s *Struct) Fields() []*Field {
	return s.fields
}
