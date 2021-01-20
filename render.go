package src

// A RenderedFile represents anything like a source code file, or xml, or json, or an image.
type RenderedFile struct {
	AbsoluteFileName string
	MimeType         string
	Buf              []byte
	Error            error
}
