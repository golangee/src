package golang

import (
	"fmt"
	"github.com/golangee/src/ast"
	"github.com/golangee/src/render"
	"strings"
)

const packageGoDocFile = "doc.go"
const mimeTypeGo = "text/x-go-source"
const mimeTypeDir = "application/x-directory"
const mimeTypeGoModule = "application/x-directory-module"

// Options for the renderer.
type Options struct {
}

// Renderer provides a go renderer.
type Renderer struct {
	opts       Options
	root       ast.Node
	importerId importerKey // track our unique key, to perform cleanup
}

// NewRenderer creates a new Renderer instance.
func NewRenderer(opts Options) *Renderer {
	return &Renderer{opts: opts}
}

// tearUp prepares the ast to be used for source generation.
func (r *Renderer) tearUp(node ast.Node) error {
	r.root = ast.Root(node)

	if err := installImporter(r); err != nil {
		return fmt.Errorf("unable to install importer: %w", err)
	}

	return nil
}

// tearDown frees allocated resources.
func (r *Renderer) tearDown() error {
	if err := uninstallImporter(r); err != nil {
		return fmt.Errorf("unable to uninstall importer: %w", err)
	}

	return nil
}

// importer resolves the current importer from the parents file.
func (r *Renderer) importer(n ast.Node) *importer {
	return importerFromTree(r, n)
}

// Render converts the given node into a render.Artifact. A partial result is returned if an error is detected.
func (r *Renderer) Render(node ast.Node) (a render.Artifact, err error) {
	if err := r.tearUp(node); err != nil {
		return nil, fmt.Errorf("unable to tearUp: %w", err)
	}

	defer func() {
		if e := r.tearDown(); e != nil && err == nil {
			err = e
		}
	}()

	root := &render.Dir{}
	err = ast.ForEachMod(node, func(mod *ast.Mod) error {
		if mod.Target.Lang == ast.LangGo {
			_, err := r.renderMod(mod, root)

			if err != nil {
				return fmt.Errorf("cannot render module '%s': %w", mod.Name, err)
			}
		}

		return nil
	})

	if err != nil {
		return root, fmt.Errorf("unable to loop modules: %w", err)
	}

	return root, nil
}

func (r *Renderer) renderMod(mod *ast.Mod, parent *render.Dir) (*render.Dir, error) {
	modDir := r.ensurePkgDir(mod.Name, parent)
	modDir.MimeType = mimeTypeGoModule

	var firstErr error

	for _, pkg := range mod.Pkgs {
		// we cannot use name here, because in go the name and the import path may be different
		if !strings.HasPrefix(pkg.Path, mod.Name+"/") {
			return nil, fmt.Errorf("declared package '%s' must be prefixed by module path '%s'", pkg.Path, mod.Name)
		}
		pkgDir := r.ensurePkgDir(pkg.Path, parent)
		files, err := r.renderPkg(pkg)
		if firstErr == nil && err != nil {
			firstErr = fmt.Errorf("unable to render package '%s': %w", pkg.Path, err)
		}

		pkgDir.Files = append(pkgDir.Files, files...)
	}

	return modDir, firstErr
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
