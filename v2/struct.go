package src

var _ NamedType = (*Struct)(nil)

// A Struct is actually a data type, like a record. Depending on the language, it can be used in a value or reference
// context. If supported, the primary use case should be the usage as a value to improve conclusiveness and
// performance by avoiding heap allocation (and potentially GC overhead). Inheritance is not possible, but other
// types may be embedded (e.g. in Go). Languages like Java use just simple classes (PoJos), because records have no
// exclusive use (they are just syntax sugar for a class with final members). In contrast to that, Go cannot express
// final fields.
type Struct struct {
	doc         string
	name        string
	visibility  Visibility
	fields      []*Field
	final       bool
	static      bool
	annotations []*Annotation
	types       []NamedType
	methods     []*Func
}

// NewStruct returns a new named struct type. A struct is always mutable, but may be used either in a value
// or pointer context. Structs are straightforward in Go but in Java just a PoJo. We do not use records, because
// they have a different semantic (read only).
func NewStruct(name string) *Struct {
	return &Struct{name: name}
}

// Static returns true, if this struct or class should pull its outer scope. This is only for Java and inner classes.
func (s *Struct) Static() bool {
	return s.static
}

// SetStatic updates the static flag. Only for Java.
func (s *Struct) SetStatic(static bool) *Struct {
	s.static = static
	return s
}

// Final returns true, if this struct or class cannot be inherited. This only applies to Java.
func (s *Struct) Final() bool {
	return s.final
}

// SetFinal updates the final flag. Only for Java.
func (s *Struct) SetFinal(final bool) *Struct {
	s.final = final
	return s
}

// Name returns the declared identifier which must be unique per package.
func (s *Struct) Name() string {
	return s.name
}

func (s *Struct) sealedNamedType() {
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

// Annotations returns the backing slice of all annotations.
func (s *Struct) Annotations() []*Annotation {
	return s.annotations
}

// AddAnnotations appends the given annotations. Note that not all render targets support type annotations, e.g.
// like Go.
func (s *Struct) AddAnnotations(a ...*Annotation) *Struct {
	s.annotations = append(s.annotations, a...)
	return s
}

// AddTypes adds a bunch of named types. This is only allowed in Java and other renderers should
// either ignore it or place them at the package level (Go).
func (s *Struct) AddTypes(types ...NamedType) *Struct {
	s.types = append(s.types, types...)
	return s
}

// Types returns the backing slice of defined named types.
func (s *Struct) Types() []NamedType {
	return s.types
}

// Methods returns all available functions.
func (s *Struct) Methods() []*Func {
	return s.methods
}

// AddMethods appends more methods to this interfaces contract.
func (s *Struct) AddMethods(f ...*Func) *Struct {
	s.methods = append(s.methods, f...)
	return s
}
