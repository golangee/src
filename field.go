package src

import "sort"

// A Field represents a (usually named) attribute or member of a struct or class.
type Field struct {
	doc         string
	visibility  Visibility
	name        string
	typeDecl    TypeDecl
	annotations []*Annotation
}

// NewField allocates a new named field. For some renderers like Go, an empty name declares an embedded type.
func NewField(name string, typeDecl TypeDecl) *Field {
	return &Field{
		name:     name,
		typeDecl: typeDecl,
	}
}

// SetVisibility updates the fields Visibility. The Go renderer will override the rendered name to match the visibility.
func (f *Field) SetVisibility(v Visibility) *Field {
	f.visibility = v
	return f
}

// Visibility returns the fields visibility.
func (f *Field) Visibility() Visibility {
	return f.visibility
}

// SetDoc updates the doc.
func (f *Field) SetDoc(doc string) *Field {
	f.doc = doc
	return f
}

// Doc returns the documentation.
func (f *Field) Doc() string {
	return f.doc
}

// SetName updates the fields name.
func (f *Field) SetName(name string) *Field {
	f.name = name
	return f
}

// Name returns the fields name.
func (f *Field) Name() string {
	return f.name
}

// AddAnnotations appends the given annotations or tags to the field.
func (f *Field) AddAnnotations(a ...*Annotation) *Field {
	f.annotations = append(f.annotations, a...)
	return f
}

// Annotations returns the backing slice of all annotations.
func (f *Field) Annotations() []*Annotation {
	return f.annotations
}

// TypeDecl returns the current type declaration.
func (f *Field) TypeDecl() TypeDecl {
	return f.typeDecl
}

// String returns a debugging representation.
func (f *Field) String() string {
	return f.name + " " + f.typeDecl.String()
}

// An Annotation represents a name and a bunch of named values. The Go renderer will emit this as a struct field
// tag. However, only the value for the empty key is used, just as is (but quoted). In Java each key represents
// the named attribute of an annotation. The name is interpreted as a fully qualified identifier.
type Annotation struct {
	name   Name
	values map[string]string
}

// NewAnnotation creates a new named Annotation. In Go the name is just interpreted as a string and has no further
// meaning.
func NewAnnotation(name Name) *Annotation {
	return &Annotation{
		name:   name,
		values: map[string]string{},
	}
}

// SetName updates the annotations name.
func (a *Annotation) SetName(name Name) *Annotation {
	a.name = name
	return a
}

// Name returns the annotations name.
func (a *Annotation) Name() Name {
	return a.name
}

// SetDefault sets the unnamed attribute value. See SetValue.
func (a *Annotation) SetDefault(value string) *Annotation {
	return a.SetValue("", value)
}

// SetValue sets a named attribute value. The value is interpreted as is, so e.g. use plain
//// values for language constants, like 3, 3.4, true or "hello world" in Java. The Go renderer only
//// ever evaluates the unnamed attribute and quotes the string itself.
func (a *Annotation) SetValue(name, value string) *Annotation {
	a.values[name] = value
	return a
}

// Value returns a specific value or the empty string.
func (a *Annotation) Value(name string) string {
	return a.values[name]
}

// Attributes returns the sorted list of attribute names.
func (a *Annotation) Attributes() []string {
	var tmp []string
	for key := range a.values {
		tmp = append(tmp, key)
	}

	sort.Strings(tmp)
	return tmp
}
