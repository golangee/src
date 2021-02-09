package ast

// A Param represents a functional input or output parameter.
type Param struct {
	ParamName        string
	ParamTypeDecl    TypeDecl
	ParamAnnotations []*Annotation
	Obj
}

// NewParam returns a new named parameter. It is valid to have unnamed parameters In go.
func NewParam(name string, typeDecl TypeDecl) *Param {
	return &Param{
		ParamName:     name,
		ParamTypeDecl: typeDecl,
	}
}

// SetName updates the current AnnotationName. An empty AnnotationName is valid for Go, but not for Java.
func (p *Param) SetName(name string) *Param {
	p.ParamName = name
	return p
}

// Name returns the parameters AnnotationName.
func (p *Param) Name() string {
	return p.ParamName
}

// SetTypeDecl updates the type declaration of the parameter.
func (p *Param) SetTypeDecl(t TypeDecl) *Param {
	p.ParamTypeDecl = t
	return p
}

// TypeDecl returns the current type declaration.
func (p *Param) TypeDecl() TypeDecl {
	return p.ParamTypeDecl
}

// String returns a debugging representation.
func (p *Param) String() string {
	return p.ParamName + " " + p.ParamTypeDecl.String()
}

// Annotations returns the backing slice of all ParamAnnotations.
func (p *Param) Annotations() []*Annotation {
	return p.ParamAnnotations
}

// AddAnnotations appends the given ParamAnnotations. Note that not all render targets support parameter ParamAnnotations, e.g.
// like Go.
func (p *Param) AddAnnotations(a ...*Annotation) *Param {
	p.ParamAnnotations = append(p.ParamAnnotations, a...)
	return p
}
