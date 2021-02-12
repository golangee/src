package golang

import (
	"fmt"
	"github.com/golangee/src/ast"
	"github.com/golangee/src/render"
	"strings"
)

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
			firstErr = fmt.Errorf("cannot render package '%s': %w", pkg.Path, err)
		}

		pkgDir.Files = append(pkgDir.Files, files...)
	}

	return modDir, firstErr
}
