package src

// A Block defines a lexical scope where introduced variables usually shadow outer variables. Depending
// on the occurrence of nesting (e.g. within a closure) other limitations (like effective final variables in Java)
// needs to be respected (not in Go, but would perhaps have been better).
type Block struct {
	elements []interface{}
}

// NewBlock creates a new lexical scope.
func NewBlock() *Block {
	return &Block{}
}

// P will render the given elements later using the %v directive.
// Special handlings have:
//  * TypeDecl will be rendered using the import logic of the according renderer.
//    This also includes the replacement of standard library types.
//  * Template contains a Go text template which is applied on the current ast node. So one has access to the entire
//    context and can generate arbitrary output.
//  * Block is recursively applied.
//  * ast.Macro is called appropriately.
func (b *Block) P(r ...interface{}) *Block {
	b.elements = append(b.elements, r...)
	return b
}

// L is like P but appends a new line at the end.
func (b *Block) L(r ...interface{}) *Block {
	b.elements = append(b.elements, r...)
	return b
}

// Invoke currently just writes a method invocation in c-style which is compatible with Java and Go. This may
// change in the future by introducing a Template.
func (b *Block) Invoke(name Name, args ...interface{}) *Block {
	b.P(name, "(")
	for i, arg := range args {
		b.P(arg)
		if i < len(args)-1 {
			b.P(", ")
		}
	}
	b.P(")")

	return b
}

// Var creates a Template to use a short declaration like the following:
//  * Java: var x = new Book()
//  * Go: x := NewBook()
func (b *Block) Var(name Name, rhs ...interface{}) *Block {
	//TODO
	return b
}

// New creates a Template to evaluate cross platform like the following:
//  * Java: new Book(a, r, g)
//  * Go: pkg.NewBook(a, r, g). This currently only works for public types. But we can pick the right one,
//    by inspecting the declarations in the current package.
func (b *Block) New(name Name, args ...interface{}) *Block {
	//TODO
	return b
}

// Call is like Invoke but uses the identifier as the target and the method name as the Name.
func (b *Block) Call(identifier, method string, args ...interface{}) *Block {
	b.P(identifier, ".")
	b.Invoke(Name(method), args...)
	return b
}

// Elements returns the backing slice of all block elements. See also P.
func (b *Block) Elements() []interface{} {
	return b.elements
}

// A Template represents a Go text template and is evaluated in its current ast context.
// The dot scope (.) represents an instance of ast.MacroCtx.
type Template string
