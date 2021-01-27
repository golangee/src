package ast2

import "strconv"

type Visibility int

func (v Visibility) String() string {
	switch v {
	case Public:
		return "public"
	case PackagePrivate:
		return "packagePrivate"
	case Protected:
		return "protected"
	case Private:
		return "private"
	default:
		return "unknown-" + strconv.Itoa(int(v))
	}
}

const (
	// Public declarations are visible for everyone.
	Public Visibility = iota

	// PackagePrivate declarations are visible only from within the package (in the sense of a module). This
	// corresponds to a Go lowercase identifier and to the Java Default rule.
	PackagePrivate

	// Protected declarations are visible only from within the package and by extending classes. This is a Java-only
	// feature. Other renderers should treat this as PackagePrivate.
	Protected

	// Private declarations are only visible within the current compilation unit or class. The semantics depends
	// on the target. This is a Java-only feature. Other renderers should treat this as PackagePrivate.
	Private
)

// Visibilities contains all available visibilities.
var Visibilities = []Visibility{Public, PackagePrivate, Protected, Private}