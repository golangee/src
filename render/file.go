package render

// A File represents anything like a source code file, or xml, or json, or an image.
type File struct {
	FileName string
	MimeType string
	Buf      []byte
	Error    error
}

func (n *File) Name() string {
	return n.FileName
}
