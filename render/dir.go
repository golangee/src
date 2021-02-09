package render

// A Dir contains other dirs and files.
type Dir struct {
	DirName string
	Files   []*File
	Dirs    []*Dir
}

func (n *Dir) Name() string {
	return n.DirName
}
