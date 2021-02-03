package ast2

// A Pkg represents a package and contains compilation units (source code files). Its position relates the local
// physical folder.
type Pkg struct {
	PkgFiles []*File
	// Path denotes the import path.
	//  * Go: the fully qualified Go path or module path for this module.
	//  * Java: the fully qualified package AnnotationName.
	Path string

	// Name denotes the actual package AnnotationName.
	//  * Go: the actual package AnnotationName, as defined by a File.
	//  * Java: the last segment (identifier) of the full qualified package AnnotationName.
	Name string

	// A Preamble comment belongs not to the actual file or package documentation and is usually
	// something like a license or generator header.
	//  * Go: first comment In a file named doc.go
	//  * Java: first comment In a file named package-info.java
	Preamble *Comment

	Obj
}

// Files appends the given files.
func (n *Pkg) Files(files ...*File) *Pkg {
	for _, file := range files {
		assertNotAttached(file)
		n.PkgFiles = append(n.PkgFiles, file)
		file.ObjParent = n
	}

	return n
}

// Children returns a defensive copy of the underlying slice. However the Node references are shared.
func (n *Pkg) Children() []Node {
	tmp := make([]Node, 0, len(n.PkgFiles)+1)

	if n.ObjComment != nil {
		tmp = append(tmp, n.Obj.ObjComment)
	}

	for _, pkg := range n.PkgFiles {
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
