package java

import (
	"github.com/golangee/src/v2"
	"sort"
)

// importer manages the rendered import section at the files top.
type importer struct {
	identifiersInScope map[string]src.Name
}

func newImporter() *importer {
	return &importer{
		identifiersInScope: map[string]src.Name{},
	}
}

// importerFromTree walks up the tree until it finds the srcFileNode. If the node is not attached to a file,
// it panics.
func importerFromTree(n node) *importer {
	root := n
	for root != nil {
		if srcNode, ok := root.(*srcFileNode); ok {
			return srcNode.importer
		}

		newRoot := root.Parent()
		if newRoot == nil {
			panic("srcFileNode not found")
		}

		root = newRoot
	}

	panic("invalid node")
}

// qualifiers returns the unique imported qualifiers.
func (p *importer) qualifiers() []string {
	tmp := map[string]string{}
	for _, name := range p.identifiersInScope {
		tmp[string(name)] = ""
	}

	var sorted []string
	for uniqueQualifier := range tmp {
		sorted = append(sorted, uniqueQualifier)
	}

	sort.Strings(sorted)

	return sorted
}

// shortify returns a qualified name, which is only valid the importers scope. It may also decide to not import
// the given name, e.g. if a collision has been detected. If the name is a universe type or not complete, the original
// name is just returned.
func (p *importer) shortify(name src.Name) src.Name {
	qual := name.Qualifier()
	id := name.Identifier()
	if id == "" || qual == "" {
		return name
	}

	otherName, inScope := p.identifiersInScope[id]
	if inScope {
		// already registered the identical qualifier, e.g.
		// a.A => A
		// a.B => B
		if otherName == name {
			return src.Name(id)
		} else {
			// name collision
			return name
		}
	}

	p.identifiersInScope[id] = name
	return src.Name(id)
}
