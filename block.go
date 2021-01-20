package src

import "fmt"

// Java contains the mimetype for Java source code files.
const Java = "text/x-java-source"

// Go contains the mimetype for Go source code files.
const Go = "text/x-go-source"

// A Block defines a lexical scope where introduced variables usually shadow outer variables. Depending
// on the occurrence of nesting (e.g. within a closure) other limitations (like effective final variables in Java)
// needs to be respected (not in Go, but would perhaps have been better).
//
// Note that the logic here is a bit weired because it can be both: at first, very low level, by using
// pure source code bytes and secondly very high level by just providing a limited set of primitives like
// Var, New, Invoke or Call and others which will be translated into the according target language.
type Block struct {
	elements []interface{}
}

// NewBlock creates a new lexical scope.
func NewBlock() *Block {
	return &Block{}
}

// P will render the given elements later using the %v directive.
// Special handlings have:
//  * TypeDecl and Name will be rendered using the import logic of the according renderer.
//    This also includes the replacement of standard library types.
//  * Block is recursively applied.
//  * ast.Macro is called appropriately.
//  * Macro is called appropriately.
// TODO why not applying any src.* type? we could do inner classes, functions etc???
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
//  * Java: var <identifier> = rhs
//  * Go: <identifier> := rhs
func (b *Block) Var(identifier string, rhs ...interface{}) *Block {
	b.P(Macro(func(ctx MacroCtx, p func(r ...interface{})) error {
		switch ctx.MimeType() {
		case Java:
			p("var ", identifier, " = ")
			p(rhs...)
		case Go:
			p(identifier, " := ")
			p(rhs...)
		default:
			return languageNotImplemented(ctx.MimeType())
		}
		return nil
	}))

	return b
}

// New creates a Template to evaluate cross platform like the following:
//  * Java: new <name>(a, r, g)
//  * Go: pkg.New<name.Identifier>(a, r, g). This currently only works for public types. But we can pick the right one,
//    by inspecting the declarations in the current package.
func (b *Block) New(name Name, args ...interface{}) *Block {
	b.P(Macro(func(ctx MacroCtx, p func(r ...interface{})) error {
		switch ctx.MimeType() {
		case Java:
			p("new ", name)
			p("(")
			for i, arg := range args {
				p(arg)
				if i > len(args)-1 {
					p(",")
				}
			}
			p(")")
		case Go:
			newName := Name(name.Qualifier() + ".New" + name.Identifier())
			p(newName)
			p("(")
			for i, arg := range args {
				p(arg)
				if i > len(args)-1 {
					p(",")
				}
			}
			p(")")
		default:
			return languageNotImplemented(ctx.MimeType())
		}
		return nil
	}))

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

// A MacroCtx allows access to actual rendering context.
type MacroCtx interface {
	// MimeType is a shortcut to the files mimetype. If not applicable, returns the empty string.
	//  * text/x-java-source
	//  * text/x-go-source
	MimeType() string

	// TODO we dont need that, if our macro contract just says that Elements are strings, Name and src.TypeDecl
	// StdLib tries to resolve the given name into a standard library type defined by the actual render context
	// without importing it.
	StdLib(name Name) Name

	// TODO we dont need that, if our macro contract just says that Elements are strings, Name and src.TypeDecl
	// Import tries to resolve the given name from the standard library and imports it for the usage in the current
	// file.
	Import(name Name) Name

	// Receiver returns the current receiver:
	//  * Java: in static contexts (like a static method) this is the empty string. In other methods just "this" is
	//    returned.
	//  * Go: in a functions context (not a struct method) this is the empty string, otherwise the receiver name
	//    name is returned. If no custom name has been defined, the first lowercase letter from the type
	//    is declared and returned.
	Receiver() string
}

// Macro is a function to be called to emit source code directly.
type Macro func(ctx MacroCtx, p func(r ...interface{})) error

func languageNotImplemented(lang string) error {
	return fmt.Errorf("language not implemented: %s", lang)
}
