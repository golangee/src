package golang

import (
	"github.com/golangee/src/ast"
	"github.com/golangee/src/render"
	"strings"
)

func (r *Renderer) renderPkg(pkg *ast.Pkg) ([]*render.File, error) {
	var res []*render.File
	var firstErr error

	if pkg.Preamble != nil || pkg.ObjComment != nil {
		tmp := &render.BufferedWriter{}
		// package license or whatever
		if pkg.Preamble != nil {
			r.writeComment(tmp, false, pkg.Name, pkg.Preamble.Text)
			tmp.Printf("\n")
		}

		// actual package comment
		if pkg.ObjComment != nil {
			r.writeComment(tmp, true, pkg.Name, pkg.ObjComment.Text)
		}

		tmp.Printf("package %s\n", pkg.Name)

		f := &render.File{
			FileName: packageGoDocFile,
			MimeType: mimeTypeGo,
			Buf:      tmp.Bytes(),
		}

		res = append(res, f)
	}

	for _, file := range pkg.PkgFiles {
		buf, err := r.renderFile(file)

		f := &render.File{
			FileName: file.Name,
			MimeType: mimeTypeGo,
		}
		f.Buf = buf
		f.Error = err

		if firstErr == nil && err != nil {
			firstErr = err
		}

		res = append(res, f)
	}

	return res, firstErr
}

// ensurePkgDir appends for each path segment a directory, if required. Returns the directory denoting
// the last segment.
func (r *Renderer) ensurePkgDir(restPath string, parent *render.Dir) *render.Dir {
	names := strings.Split(restPath, "/")

	dir := parent.Directory(names[0])
	if dir == nil {
		dir = &render.Dir{DirName: names[0], MimeType: mimeTypeDir}
		parent.Dirs = append(parent.Dirs, dir)
	}

	if len(names) == 1 {
		return dir
	}

	return r.ensurePkgDir(strings.Join(names[1:], "/"), dir)
}


