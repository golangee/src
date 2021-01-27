package ast2

// File represents a physical source code file respective compilation unit.
//  * Go: <lowercase name>.go
//  * Java: <CamelCasePrimaryTypeName>.java
type File struct {
	// A Preamble comment belongs not to any type and is usually
	// something like a license or generator header as the first comment in the actual file.
	// The files comment is actually Obj.Comment.
	Preamble *Comment
	Types    []Node
	Funcs    []Node
	Obj
}
