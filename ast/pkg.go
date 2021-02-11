package ast

import "strings"

// A Pkg represents a package and contains compilation units (source code files). Its position relates the local
// physical folder.
type Pkg struct {
	PkgFiles []*File
	// Path denotes the import path (separated with slashes, even on windows).
	//  * Go: the fully qualified Go path or module path for this module.
	//  * Java: the fully qualified package name.
	Path string

	// Name denotes the actual package name.
	//  * Go: the actual package name, as defined by a File.
	//  * Java: the last segment (identifier) of the full qualified package name.
	Name string

	// A Preamble comment belongs not to the actual file or package documentation and is usually
	// something like a license or generator header.
	//  * Go: first comment In a file named doc.go
	//  * Java: first comment In a file named package-info.java
	Preamble *Comment

	Obj
}

// NewPkg creates a new package with the given Path and sets the Name to the last segment. Path names must be
// separated by / (even on Windows).
func NewPkg(path string) *Pkg {
	names := strings.Split(path, "/")
	return &Pkg{
		Path: path,
		Name: names[len(names)-1],
	}
}

// SetPreamble sets a non-package comment.
func (n *Pkg) SetPreamble(text string) *Pkg {
	n.Preamble = NewComment(text)
	return n
}

// SetComment sets the package comment section.
func (n *Pkg) SetComment(text string) *Pkg {
	n.ObjComment = NewComment(text)
	return n
}

// SetName updates the name. Some language targets may ignore that, like Java.
func (n *Pkg) SetName(name string) *Pkg {
	n.Name = name
	return n
}

// Files appends the given files.
func (n *Pkg) AddFiles(files ...*File) *Pkg {
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
