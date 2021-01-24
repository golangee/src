package ast2

// A Pkg represents a package and contains compilation units (source code files).
type Pkg struct {
	PkgFiles    []*File
	PkgPackages []*Pkg
	// Path denotes the import path.
	//  * Go: the fully qualified Go path or module path for this module.
	//  * Java: the fully qualified package name.
	Path string

	// Name denotes the actual package name.
	//  * Go: the actual package name, as defined by a File.
	//  * Java: the last segment (identifier) of the full qualified package name.
	Name string

	// A Preamble comment belongs not to the actual file or package documentation and is usually
	// something like a license or generator header.
	//  * Go: first comment in a file named doc.go
	//  * Java: first comment in a file named package-info.java
	Preamble *Comment

	Obj
}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (n *Pkg) Children() []Node {
	tmp := make([]Node, 0, len(n.PkgFiles)+len(n.PkgPackages)+1)
	tmp = append(tmp, n.Obj.ObjComment)

	for _, pkg := range n.PkgFiles {
		tmp = append(tmp, pkg)
	}

	for _, pkg := range n.PkgPackages {
		tmp = append(tmp, pkg)
	}

	return tmp
}

// Doc sets the nodes comment.
func (n *Pkg) Doc(text string) *Pkg {
	n.Obj.ObjComment = NewComment(text)
	n.Obj.ObjComment.ObjParent = n
	return n
}
