package ast2

// A Mod is the root of a project and describes a module with packages.
//  * Java: denotes a gradle module.
//  * Go: describes a Go module.
type Mod struct {
	Pkgs []*Pkg
	Obj
}

// NewModule allocates a new Module.
func NewModule() *Mod {
	return &Mod{}
}

// Packages appends the given packages and updates the Parent accordingly.
func (n *Mod) Packages(packages ...*Pkg) *Mod {
	n.Pkgs = append(n.Pkgs, packages...)
	for _, pkg := range packages {
		pkg.Obj.ObjParent = n
	}

	return n
}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (n *Mod) Children() []Node {
	tmp := make([]Node, 0, len(n.Pkgs)+1)
	tmp = append(tmp, n.Obj.ObjComment)

	for _, pkg := range n.Pkgs {
		tmp = append(tmp, pkg)
	}

	return tmp
}

// Doc sets the nodes comment.
func (n *Mod) Doc(text string) *Mod {
	n.Obj.ObjComment = NewComment(text)
	n.Obj.ObjComment.ObjParent = n
	return n
}
