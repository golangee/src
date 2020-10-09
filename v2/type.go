package src

// NamedType declares a new named type which is either derived from another existing type
// or is a struct or interface.
type NamedType interface {
	Name() string
	isNamedType()
}
