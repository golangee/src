package golang

import (
	"fmt"
	"github.com/golangee/src/ast"
	"github.com/golangee/src/render"
)

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
		return fmt.Errorf("unable to uninstall importer: %w")
	}

	return nil
}

// importer resolves the current importer from the parents file.
func (r *Renderer) importer(n ast.Node) *importer {
	return importerFromTree(r, n)
}

// Render converts the given node into a render.Artifact.
func (r *Renderer) Render(node ast.Node) (a render.Artifact, err error) {
	if err := r.tearUp(node); err != nil {
		return nil, fmt.Errorf("unable to tearUp: %w", err)
	}

	defer func() {
		if e := r.tearDown(); e != nil && err == nil {
			err = e
		}
	}()

	return nil, nil
}
