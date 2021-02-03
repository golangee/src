package ast2

var _ NamedType = (*Struct)(nil)

// A Struct is actually a data type, like a record. Depending on the language, it can be used in a value or reference
// context. If supported, the primary use case should be the usage as a value to improve conclusiveness and
// performance by avoiding heap allocation (and potentially GC overhead). Inheritance is not possible, but other
// types may be embedded (e.g. in Go). Languages like Java use just simple classes (PoJos), because records have no
// exclusive use (they are just syntax sugar for a class with final members). In contrast to that, Go cannot express
// final fields.
type Struct struct {
	TypeName        string
	TypeVisibility  Visibility
	TypeFields      []*Field
	TypeStatic      bool
	TypeAnnotations []*Annotation
	TypeMethods     []*Func
	Types           []NamedType // only valid for language which can declare named nested type like java
	Obj
}

// NewStruct returns a new named struct type. A struct is always mutable, but may be used either in a value
// or pointer context. Structs are straightforward in Go but in Java just a PoJo. We do not use records, because
// they have a different semantic (read only).
func NewStruct(name string) *Struct {
	return &Struct{TypeName: name}
}

// Static returns true, if this struct or class should pull its outer scope. This is only for Java and inner classes.
func (s *Struct) Static() bool {
	return s.TypeStatic
}

// SetStatic updates the static flag. Only for Java.
func (s *Struct) SetStatic(static bool) *Struct {
	s.TypeStatic = static
	return s
}

// Name returns the declared identifier which must be unique per package.
func (s *Struct) Name() string {
	return s.TypeName
}

func (s *Struct) sealedNamedType() {
	panic("implement me")
}

// SetVisibility sets the visibility. The default is Public.
func (s *Struct) SetVisibility(v Visibility) *Struct {
	s.TypeVisibility = v
	return s
}

// Visibility returns the current visibility. The default is Public.
func (s *Struct) Visibility() Visibility {
	return s.TypeVisibility
}

// AddFields appends the given fields to the struct.
func (s *Struct) AddFields(fields ...*Field) *Struct {
	for _, field := range fields {
		assertNotAttached(field)
		assertSettableParent(field).SetParent(field)
		s.TypeFields = append(s.TypeFields, field)
	}
	return s
}

// Fields returns the currently configured fields.
func (s *Struct) Fields() []*Field {
	return s.TypeFields
}

// Annotations returns the backing slice of all annotations.
func (s *Struct) Annotations() []*Annotation {
	return s.TypeAnnotations
}

// AddAnnotations appends the given annotations. Note that not all render targets support type annotations, e.g.
// like Go.
func (s *Struct) AddAnnotations(a ...*Annotation) *Struct {
	for _, annotation := range a {
		assertNotAttached(annotation)
		assertSettableParent(annotation).SetParent(s)
	}

	return s
}

// Methods returns all available functions.
func (s *Struct) Methods() []*Func {
	return s.TypeMethods
}

// AddMethods appends more methods to this interfaces contract.
func (s *Struct) AddMethods(f ...*Func) *Struct {
	for _, fun := range f {
		assertNotAttached(fun)
		assertSettableParent(fun).SetParent(s)
		s.TypeMethods = append(s.TypeMethods, fun)
	}

	return s
}

func (s *Struct) NamedTypes() []NamedType {
	return s.Types
}

func (s *Struct) AddNamedTypes(t ...NamedType) *Struct {
	for _, namedType := range t {
		assertNotAttached(namedType)
		assertSettableParent(namedType).SetParent(s)
		s.Types = append(s.Types, namedType)
	}

	return s
}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (s *Struct) Children() []Node {
	tmp := make([]Node, 0, len(s.TypeFields)+len(s.TypeAnnotations)+len(s.TypeMethods)+len(s.Types))
	for _, param := range s.TypeAnnotations {
		tmp = append(tmp, param)
	}

	for _, param := range s.TypeFields {
		tmp = append(tmp, param)
	}

	for _, param := range s.TypeMethods {
		tmp = append(tmp, param)
	}

	for _, namedType := range s.Types {
		tmp = append(tmp, namedType)
	}

	return tmp
}
