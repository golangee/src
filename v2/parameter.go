package src

// A Param represents a functional input or output parameter.
type Param struct {
	doc      string
	name     string
	typeDecl TypeDecl
}

// NewParam returns a new named parameter. It is valid to have unnamed parameters in go.
func NewParam(name string, typeDecl TypeDecl) *Param {
	return &Param{
		name:     name,
		typeDecl: typeDecl,
	}
}

// SetDoc updates the parameters comment. Go does not have an explicit representation,
// however the text is just appended to the functions comment. The best is to use
// the ellipsis. In Java this is also merged into methods comment but with using the @param annotation.
func (p *Param) SetDoc(doc string) *Param {
	p.doc = doc
	return p
}

// Doc returns the current comment.
func (p *Param) Doc() string {
	return p.doc
}

// SetName updates the current name. An empty name is valid for Go, but not for Java.
func (p *Param) SetName(name string) *Param {
	p.name = name
	return p
}

// SetTypeDecl updates the type declaration of the parameter.
func (p *Param) SetTypeDecl(t TypeDecl) *Param {
	p.typeDecl = t
	return p
}

// TypeDecl returns the current type declaration.
func (p *Param) TypeDecl() TypeDecl {
	return p.typeDecl
}

// String returns a debugging representation.
func (p *Param) String() string {
	return p.name + " " + p.typeDecl.String()
}
