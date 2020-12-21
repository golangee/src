package golang

import (
	"github.com/golangee/src/v2"
	"github.com/golangee/src/v2/ast"
	"sort"
	"strconv"
	"strings"
)

type importerKey int

const importerId importerKey = 1

// importer manages the rendered import section at the files top.
type importer struct {
	namedImports map[string]string // named import => qualifier
}

func newImporter() *importer {
	return &importer{
		namedImports: map[string]string{},
	}
}

// installImporter installs a new importer instance into every ast.SrcFileNode.
func installImporter(n *ast.ModNode) {
	for _, node := range n.Packages() {
		for _, fileNode := range node.Files() {
			fileNode.SetValue(importerId, newImporter())
		}
	}
}

// importerFromTree walks up the tree until it finds the first importer from any ast.Node.Value.
func importerFromTree(n ast.Node) *importer {
	root := n
	for root != nil {
		if imp, ok := root.Value(importerId).(*importer); ok {
			return imp
		}

		newRoot := root.Parent()
		if newRoot == nil {
			panic("no attached importer found in ast scope")
		}

		root = newRoot
	}

	panic("invalid node")
}

// qualifiers returns the unique imported qualifiers.
func (p *importer) qualifiers() []string {
	tmp := map[string]string{}
	for _, name := range p.namedImports {
		tmp[string(name)] = ""
	}

	var sorted []string
	for uniqueQualifier := range tmp {
		sorted = append(sorted, uniqueQualifier)
	}

	sort.Strings(sorted)

	return sorted
}

// shortify returns a qualified name, which is only valid in the importers scope. It may also decide to not import
// the given name, e.g. if a collision has been detected. If the name is a universe type or not complete, the original
// name is just returned.
func (p *importer) shortify(name src.Name) src.Name {
	qual := name.Qualifier()
	id := name.Identifier()
	if id == "" || qual == "" {
		return name
	}

	namedImportIdxName := strings.LastIndex(qual, "/") // e.g. 3 for net/http or -1 for net
	if namedImportIdxName == -1 {
		namedImportIdxName = 0
	}

	namedImport := MakePrivate(MakeIdentifier(qual[namedImportIdxName:]))

	otherQualifier, inScope := p.namedImports[namedImport]
	if inScope {
		// already registered the identical qualifier, e.g.
		// net/http => http
		if otherQualifier == qual {
			return src.Name(namedImport + "." + id)
		} else {
			// name collision, build something artificial with increasing number
			num := 1
			for {
				num++
				namedImport2 := namedImport + strconv.Itoa(num)
				otherQualifier, inScope = p.namedImports[namedImport]
				if inScope {
					if otherQualifier == qual {
						return src.Name(namedImport + "." + id)
					}
					// loop again until either found or no other entry found
				} else {
					namedImport = namedImport2
					break
				}
			}
		}
	}

	p.namedImports[namedImport] = qual
	return src.Name(namedImport + "." + id)
}
