package golang

import "github.com/golangee/src/v2"

type modNode struct {
	module   *src.Module
	packages []*pkgNode
}

type pkgNode struct {
	parent   *modNode
	pkg      *src.Package
	srcFiles []*srcFileNode
}

type srcFileNode struct {
	parent *pkgNode
	file   *src.SrcFile
	types  []*typeNode
}

type typeNode struct {
	parent    *srcFileNode
	namedType src.NamedType
}

func newModNode(mod *src.Module) *modNode {
	n := &modNode{
		module: mod,
	}

	for _, p := range mod.Packages() {
		n.packages = append(n.packages, newPkgNode(n, p))
	}

	return n
}

func newPkgNode(parent *modNode, pkg *src.Package) *pkgNode {
	n := &pkgNode{
		parent: parent,
		pkg:    pkg,
	}

	for _, file := range pkg.SrcFiles() {
		n.srcFiles = append(n.srcFiles, newSrcFileNode(n, file))
	}

	return n
}

func newSrcFileNode(parent *pkgNode, file *src.SrcFile) *srcFileNode {
	n := &srcFileNode{
		parent: parent,
		file:   file,
	}

	for _, namedType := range file.Types() {
		n.types = append(n.types, newTypeNode(n, namedType))
	}

	return n
}

func newTypeNode(parent *srcFileNode, t src.NamedType) *typeNode {
	return &typeNode{
		parent:    parent,
		namedType: t,
	}
}
