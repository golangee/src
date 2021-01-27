// Copyright 2020 Torben Schinke
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ast2

import (
	"github.com/golangee/src/stdlib"
	"strconv"
)

// A TypeDecl is the most complex part to represent. Examples of valid type declarations are:
//  * Go/Java: int => SimpleTypeDecl
//  * Go: *int => TypeDeclPtr
//  * Go: map[string]int or Java: java.util.Map<String,Integer> => GenericTypeDecl with name "map"
//  * Go: []int or Java: int[] => SliceTypeDecl
//  * Go: [5]int => ArrayTypeDecl
//  * Go draft: List[T int] or Java: List<Integer> => GenericTypeDecl
//  * Go: func(a, b int) (string, error) => FuncTypeDecl
//  * Java: X <T extends List,V extends List<? super Integer>> => GenericTypeDecl + NamedTypeDecl
//  * Go: chan<- string or <-chan string or chan string => ChanTypeDecl
type TypeDecl interface {
	// String returns a human readable declaration for debugging purposes.
	String() string

	// sealedTypeDecl ensures that nobody can implement his own TypeDecl which is generally not acceptable for
	// our code generators.
	sealedTypeDecl()

	Node
}

// TypeDeclPtr is a TypeDecl declares a pointer type.
type TypeDeclPtr struct {
	Decl TypeDecl
	Obj
}

// NewTypeDeclPtr returns a new wrapper around another type declaration. This is only possible on a language
// level for Go but for Java java.util.concurrent.atomic.AtomicReference is used.
func NewTypeDeclPtr(decl TypeDecl) *TypeDeclPtr {
	t := &TypeDeclPtr{Decl: decl}
	assertNotAttached(decl)
	assertSettableParent(decl).SetParent(t)

	return t
}

// TypeDecl returns the referenced type.
func (t *TypeDeclPtr) TypeDecl() TypeDecl {
	return t.Decl
}

// String returns a debugging representation.
func (t *TypeDeclPtr) String() string {
	return "*" + t.Decl.String()
}

func (t *TypeDeclPtr) sealedTypeDecl() {
	panic("sealed type")
}

//======

// A SimpleTypeDecl just refers to a named type.
type SimpleTypeDecl struct {
	SimpleName Name
	Obj
}

// NewSimpleTypeDecl create a simple type declaration.
func NewSimpleTypeDecl(name Name) *SimpleTypeDecl {
	return &SimpleTypeDecl{SimpleName: name}
}

// Name is the qualifier and identifier of the type.
func (t *SimpleTypeDecl) Name() Name {
	return t.SimpleName
}

// SetName updates the qualified name.
func (t *SimpleTypeDecl) SetName(name Name) {
	t.SimpleName = name
}

// String returns a debugging representation.
func (t *SimpleTypeDecl) String() string {
	return string(t.SimpleName)
}

func (t *SimpleTypeDecl) sealedTypeDecl() {
	panic("sealed type")
}

//======

// A GenericTypeDecl refers to a named type and contains (optional) named type parameters, commonly known as generics.
type GenericTypeDecl struct {
	TypeDecl   TypeDecl
	TypeParams []TypeDecl
	Obj
}

// NewGenericDecl creates and returns a new generic declaration.
func NewGenericDecl(typeDecl TypeDecl, params ...TypeDecl) *GenericTypeDecl {
	t := &GenericTypeDecl{
		TypeDecl:   typeDecl,
		TypeParams: params,
	}

	assertNotAttached(typeDecl)
	assertSettableParent(typeDecl).SetParent(t)

	for _, param := range params {
		assertNotAttached(param)
		assertSettableParent(param).SetParent(t)
	}

	return t
}

// Type returns the actual Type definition.
func (t *GenericTypeDecl) Type() TypeDecl {
	return t.TypeDecl
}

// SetType updates the type declaration.
func (t *GenericTypeDecl) SetType(typeDecl TypeDecl) {
	t.TypeDecl = typeDecl
}

// Params returns all declared type parameters.
func (t *GenericTypeDecl) Params() []TypeDecl {
	return t.TypeParams
}

// String returns a debugging representation.
func (t *GenericTypeDecl) String() string {
	tmp := t.TypeDecl.String()
	tmp += "<"
	for i, param := range t.TypeParams {
		tmp += param.String()
		if i < len(t.TypeParams)-1 {
			tmp += ","
		}
	}
	tmp += ">"

	return tmp
}

func (t *GenericTypeDecl) sealedTypeDecl() {
	panic("sealed type")
}

//======

// A NamedTypeDecl tupels a normal string identifier, usually something abstract like T and a concrete type. Note
// that ? may have a special meaning, e.g. for Java, where it declares an anonymous wildcard. It may also specify a
// bound.
type NamedTypeDecl struct {
	TypeName  string
	TypeBound TypeBound
	TypeDecl  TypeDecl
	Obj
}

// NewNamedTypeDecl returns
func NewNamedTypeDecl(name string, decl TypeDecl) *NamedTypeDecl {
	t := &NamedTypeDecl{
		TypeName:  name,
		TypeDecl:  decl,
		TypeBound: UnboundedType,
	}

	assertNotAttached(decl)
	assertSettableParent(decl).SetParent(t)

	return t
}

// Bounds returns either UnboundedType or UpperBoundedType or LowerBoundedType.
func (t *NamedTypeDecl) Bound() TypeBound {
	return t.TypeBound
}

// SetBound updates the bound.
func (t *NamedTypeDecl) SetBound(bound TypeBound) {
	t.TypeBound = bound
}

// Type returns the actual Type definition.
func (t *NamedTypeDecl) Type() TypeDecl {
	return t.TypeDecl
}

// SetType updates the type declaration.
func (t *NamedTypeDecl) SetType(typeDecl TypeDecl) {
	t.TypeDecl = typeDecl
}

// Name returns the introduced parameter name. Take care of ? for special names, like wildcards.
func (t *NamedTypeDecl) Name() string {
	return t.TypeName
}

// SetName updates the named type declaration.
func (t *NamedTypeDecl) SetName(name string) {
	t.TypeName = name
}

// String returns a debugging representation.
func (t *NamedTypeDecl) String() string {
	if t.TypeBound == UnboundedType {
		return t.TypeName + " " + t.TypeDecl.String()
	}

	return t.TypeName + " " + string(t.TypeBound) + " " + t.TypeDecl.String()
}

func (t *NamedTypeDecl) sealedTypeDecl() {
	panic("sealed type")
}

// A TypeBound is currently only used by the Java renderer and is used to declare upper or lower type bounds.
type TypeBound string

const (
	// UnboundedType declares no special bound, which is the default. Remember the PECS rule (producer extends,
	// consumer super).
	UnboundedType TypeBound = ""
	// UpperBoundedType declares to allow the named type or its sub-types. E.g. List<? extends Number> allows Number
	// and sub-types likes Integer.
	UpperBoundedType TypeBound = "extends"
	// LowerBoundedType declares to allow the named type or any of its super-types. E.g. List<? super Integer> allows
	// Integer or any of its super-types likes Number.
	LowerBoundedType TypeBound = "super"
)

//======

// A SliceTypeDecl declares a slice of an arbitrary type. This is also used for Java arrays, because Java cannot
// express a fixed length array type.
type SliceTypeDecl struct {
	TypeDecl TypeDecl
	Obj
}

// NewSliceTypeDecl returns a new slice type declaration.
func NewSliceTypeDecl(typeDecl TypeDecl) *SliceTypeDecl {
	t := &SliceTypeDecl{TypeDecl: typeDecl}
	assertNotAttached(typeDecl)
	assertSettableParent(typeDecl).SetParent(t)

	return t
}

// Type returns the declared type.
func (t *SliceTypeDecl) Type() TypeDecl {
	return t.TypeDecl
}

// SetType updates the named type declaration.
func (t *SliceTypeDecl) SetType(typeDecl TypeDecl) *SliceTypeDecl {
	t.TypeDecl = typeDecl
	return t
}

// String returns a debugging representation.
func (t *SliceTypeDecl) String() string {
	return "[]" + t.TypeDecl.String()
}

func (t *SliceTypeDecl) sealedTypeDecl() {
	panic("sealed type")
}

//====== //TODO continue

// ArrayTypeDecl declares a fixed length array type of an arbitrary type. This is not expressible in Java and
// degenerates to a normal array.
type ArrayTypeDecl struct {
	len      int
	typeDecl TypeDecl
}

// NewArrayTypeDecl returns a new array type.
func NewArrayTypeDecl(len int, typeDecl TypeDecl) *ArrayTypeDecl {
	return &ArrayTypeDecl{
		len:      len,
		typeDecl: typeDecl,
	}
}

// TypeDecl returns the declared type.
func (t *ArrayTypeDecl) TypeDecl() TypeDecl {
	return t.typeDecl
}

// SetTypeDecl updates the named type declaration.
func (t *ArrayTypeDecl) SetTypeDecl(typeDecl TypeDecl) *ArrayTypeDecl {
	t.typeDecl = typeDecl
	return t
}

// Len returns the declared array length.
func (t *ArrayTypeDecl) Len() int {
	return t.len
}

// String returns a debugging representation.
func (t *ArrayTypeDecl) String() string {
	return "[" + strconv.Itoa(t.len) + "]" + t.typeDecl.String()
}

func (t *ArrayTypeDecl) sealedTypeDecl() {
	panic("sealed type")
}

//======

// NewMapDecl is just a normal generic declaration but represents either the Go builtin type "map" or
// for Java the java.util.Map type.
func NewMapDecl(key, val TypeDecl) *GenericTypeDecl {
	return NewGenericDecl(NewSimpleTypeDecl(stdlib.Map), key, val)
}

//======

// NewListDecl is just a normal generic declaration but represents either the Go slice type or
// for Java the java.util.List type.
func NewListDecl(val TypeDecl) *GenericTypeDecl {
	return NewGenericDecl(NewSimpleTypeDecl(stdlib.List), val)
}

//======

// ChanDir indicates the channel direction.
type ChanDir string

const (
	// ChanSendRecv indicates a bidirectional go channel.
	ChanSendRecv ChanDir = "chan"
	// ChanSend represents a sending-only channel declaration.
	ChanSend ChanDir = "chan<-"
	// ChanRecv represents a receiving-only channel declaration.
	ChanRecv ChanDir = "<-chan"
)

// ChanTypeDecl declares a sending or receiving or a bidirectional channel. In Go, this is part of the
// language itself. Java does not have such a type, but java.util.concurrent.BlockingQueue may be the
// equivalent.
type ChanTypeDecl struct {
	typeDecl TypeDecl
	dir      ChanDir
}

// NewChanTypeDecl creates a bidirectional channel type.
func NewChanTypeDecl(typeDecl TypeDecl) *ChanTypeDecl {
	return &ChanTypeDecl{
		typeDecl: typeDecl,
		dir:      ChanSendRecv,
	}
}

// TypeDecl returns the declared type.
func (t *ChanTypeDecl) TypeDecl() TypeDecl {
	return t.typeDecl
}

// SetTypeDecl updates the named type declaration.
func (t *ChanTypeDecl) SetTypeDecl(typeDecl TypeDecl) *ChanTypeDecl {
	t.typeDecl = typeDecl
	return t
}

// String returns a debugging representation.
func (t *ChanTypeDecl) String() string {
	return string(t.dir) + " " + t.typeDecl.String()
}

func (t *ChanTypeDecl) sealedTypeDecl() {
	panic("sealed type")
}

//======

// A FuncTypeDecl is only valid for Go, because in Java this is not directly expressible, and requires a
// "functional interface" which would be just a SimpleTypeDecl.
type FuncTypeDecl struct {
	in  []*Param
	out []*Param
}

// NewFuncTypeDecl returns a parameterless function signature.
func NewFuncTypeDecl() *FuncTypeDecl {
	return &FuncTypeDecl{}
}

// AddInputParams appends the given parameters as the functions input parameters.
func (f *FuncTypeDecl) AddInputParams(p ...*Param) *FuncTypeDecl {
	f.in = append(f.in, p...)
	return f
}

// InputParams returns the current functions input parameters.
func (f *FuncTypeDecl) InputParams() []*Param {
	return f.in
}

// AddOutputParams appends the given parameters as the functions output parameters.
func (f *FuncTypeDecl) AddOutputParams(p ...*Param) *FuncTypeDecl {
	f.out = append(f.out, p...)
	return f
}

// OutputParams returns the current functions output parameters.
func (f *FuncTypeDecl) OutputParams() []*Param {
	return f.out
}

// String returns a debugging representation.
func (f *FuncTypeDecl) String() string {
	tmp := "func("
	for i, param := range f.in {
		tmp += param.String()
		if i < len(f.in)-1 {
			tmp += ","
		}
	}
	tmp += ")"

	if len(f.out) > 0 {
		tmp += " "
	}

	if len(f.out) > 1 {
		tmp += "("
	}

	for i, param := range f.out {
		tmp += param.String()
		if i < len(f.in)-1 {
			tmp += ","
		}
	}

	if len(f.out) > 1 {
		tmp += ")"
	}

	return tmp
}

func (f *FuncTypeDecl) sealedTypeDecl() {
	panic("sealed type")
}
