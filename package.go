package src

import "strings"

// A Package is part of a module and contains only files. Packages cannot be nested.
type Package struct {
	importPath  string
	packageName string
	doc         string
	docPreamble string
	files       []*SrcFile
}

// NewPackage is the factory to create an empty package. The import path may be different from the packageName.
// This depends solely on the renderer. For the java.Render and golang.Render the importPath denotes the physical
// folder structure and packageName the package name within a file. Note that this excludes different package
// names within a folder, which is generally possible in go (at least fully allowed for *_test blackbox tests).
// In java the packageName must be equivalent to the last importPath segment. The importPath must be always a slash
// / separated path.
func NewPackage(importPath, packageName string) *Package {
	return &Package{importPath: importPath, packageName: packageName}
}

// AddSrcFiles appends other files into this package.
func (p *Package) AddSrcFiles(files ...*SrcFile) *Package {
	p.files = append(p.files, files...)
	return p
}

// SrcFiles returns the backing slice of the contained files.
func (p *Package) SrcFiles() []*SrcFile {
	return p.files
}

// ImportPath returns the according path.
func (p *Package) ImportPath() string {
	return p.importPath
}

// PackageName returns the according name. If empty, it returns the last segment of the ImportPath.
func (p *Package) PackageName() string {
	if p.packageName == "" && len(p.importPath) > 0 {
		names := strings.Split(p.importPath, "/")
		if len(names) > 0 {
			return names[len(names)-1]
		}
	}
	return p.packageName
}

// SetDocPreamble updates the preamble section in the file, which is not considered a package documentation.
func (p *Package) SetDocPreamble(doc string) *Package {
	p.docPreamble = doc
	return p
}

// DocPreamble returns the preamble section in the file, which is not considered a package documentation.
func (p *Package) DocPreamble() string {
	return p.docPreamble
}

// SetDoc sets the package documentation, which is e.g. emitted to a doc.go or a package-info.java.
func (p *Package) SetDoc(doc string) *Package {
	p.doc = doc
	return p
}

// Doc returns the package documentation.
func (p *Package) Doc() string {
	return p.doc
}
