package src

// A Module contains multiple packages and files.
type Module struct {
	packages []*Package
}

// NewModule allocates a new Module.
func NewModule() *Module {
	return &Module{}
}

// AddPackages appends the given packages.
func (m *Module) AddPackages(packages ...*Package) *Module {
	m.packages = append(m.packages, packages...)
	return m
}

// Packages returns the backing slice for the packages.
func (m *Module) Packages() []*Package {
	return m.packages
}
