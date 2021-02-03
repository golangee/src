package ast2

// A Field represents a (usually named) attribute or member of a struct or class.
type Field struct {
	FieldVisibility  Visibility
	FieldName        string
	FieldType        TypeDecl
	FieldAnnotations []*Annotation
	Obj
}

// NewField allocates a new named field. For some renderers like Go, an empty name declares an embedded type.
func NewField(name string, typeDecl TypeDecl) *Field {
	return &Field{
		FieldName: name,
		FieldType: typeDecl,
	}
}

// SetVisibility updates the fields Visibility. The Go renderer will override the rendered name to match the visibility.
func (f *Field) SetVisibility(v Visibility) *Field {
	f.FieldVisibility = v
	return f
}

// Visibility returns the fields visibility.
func (f *Field) Visibility() Visibility {
	return f.FieldVisibility
}

// SetName updates the fields name.
func (f *Field) SetName(name string) *Field {
	f.FieldName = name
	return f
}

// Name returns the fields name.
func (f *Field) Name() string {
	return f.FieldName
}

// AddAnnotations appends the given annotations or tags to the field.
func (f *Field) AddAnnotations(a ...*Annotation) *Field {
	f.FieldAnnotations = append(f.FieldAnnotations, a...)
	return f
}

// Annotations returns the backing slice of all annotations.
func (f *Field) Annotations() []*Annotation {
	return f.FieldAnnotations
}

// TypeDecl returns the current type declaration.
func (f *Field) TypeDecl() TypeDecl {
	return f.FieldType
}

// String returns a debugging representation.
func (f *Field) String() string {
	return f.FieldName + " " + f.FieldType.String()
}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (f *Field) Children() []Node {
	tmp := make([]Node, 0, len(f.FieldAnnotations)+1)
	tmp = append(tmp, f.FieldType)
	for _, annotation := range f.FieldAnnotations {
		tmp = append(tmp, annotation)
	}

	return tmp
}
