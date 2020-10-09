package src

// SrcFile represents a physical source code file.
type SrcFile struct {
	name        string
	types       []NamedType
	docPreamble string
	doc         string
}

// NewSrcFile creates a new source file. Do not add the according file extension. E.g. instead of using
// MyClass.java just use MyClass or instead of mystruct.go use just mystruct.
func NewSrcFile(name string) *SrcFile {
	return &SrcFile{name: name}
}

// Name returns the file name of this file.
func (f *SrcFile) Name() string {
	return f.name
}

// AddTypes adds a bunch of named types.
func (f *SrcFile) AddTypes(types ...NamedType) *SrcFile {
	f.types = append(f.types, types...)
	return f
}

// Types returns the backing slice of defined named types.
func (f *SrcFile) Types() []NamedType {
	return f.types
}

// SetDocPreamble updates the preamble section in the file, which is not considered a package documentation.
func (f *SrcFile) SetDocPreamble(doc string) *SrcFile {
	f.docPreamble = doc
	return f
}

// DocPreamble returns the preamble section in the file, which is not considered a package documentation.
func (f *SrcFile) DocPreamble() string {
	return f.docPreamble
}

// SetDoc updates the package documentation.
func (f *SrcFile) SetDoc(doc string) *SrcFile {
	f.doc = doc
	return f
}

// Doc returns the package documentation.
func (f *SrcFile) Doc() string {
	return f.doc
}
