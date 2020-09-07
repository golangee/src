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

package src

import "strconv"

// Types and generics are expressed by a type declaration. For example:
//   int: TypeDecl{qualifier:"int"}
//     is equal to NewTypeDecl("int")
//   []int: TypeDecl{qualifier:"[]",params:[]TypeDecl{qualifier:"int"}}
//     is equal to NewSliceDecl("int")
type TypeDecl struct {
	qualifier                 Qualifier
	params                    []*TypeDecl
	parent                    FileProvider
	funcParams, funcRetParams []*Parameter
}

func NewCallDecl(qualifier Qualifier) *TypeDecl {
	return &TypeDecl{
		qualifier: qualifier,
	}
}

func NewTypeDecl(qualifier Qualifier) *TypeDecl {
	return &TypeDecl{qualifier: qualifier}
}

func NewGenericDecl(qualifier Qualifier, params ...*TypeDecl) *TypeDecl {
	return &TypeDecl{
		qualifier: qualifier,
		params:    params,
	}
}

func NewSliceDecl(t *TypeDecl) *TypeDecl {
	return NewGenericDecl("[]", t)
}

func NewArrayDecl(len int64, t *TypeDecl) *TypeDecl {
	return NewGenericDecl(Qualifier("["+strconv.FormatInt(len, 10)+"]"), t)
}

func NewChanDecl(t *TypeDecl) *TypeDecl {
	return NewGenericDecl("chan", t)
}

func NewMapDecl(key, val *TypeDecl) *TypeDecl {
	return NewGenericDecl("map", key, val)
}

func NewPointerDecl(t *TypeDecl) *TypeDecl {
	return NewGenericDecl("*", t)
}

func NewErrorDecl() *TypeDecl {
	return NewTypeDecl("error")
}

func NewFuncDecl(params []*Parameter, returns []*Parameter) *TypeDecl {
	d := NewGenericDecl("func")
	d.funcParams = params
	d.funcRetParams = returns
	return d
}

// Clones performs a deep copy but without a parent
func (t *TypeDecl) Clone() *TypeDecl {
	c := &TypeDecl{
		qualifier: t.qualifier,
	}

	for _, param := range t.params {
		c.params = append(c.params, param.Clone())
	}

	return c
}

func (t *TypeDecl) onAttach(parent FileProvider) {
	if t == nil {
		return
	}

	t.parent = parent
	for _, p := range t.params {
		p.onAttach(parent)
	}
}

// Qualifier of the type declaration.
func (t *TypeDecl) Qualifier() Qualifier {
	return t.qualifier
}

// Params are the generic parameters.
func (t *TypeDecl) Params() []*TypeDecl {
	return t.params
}

func (t *TypeDecl) Emit(w Writer) {
	base := t.parent.File().Use(t.qualifier)
	w.Printf(base)
	w.Printf(" ")
	switch base {
	case "map":
		w.Printf("[")
		t.params[0].Emit(w)
		w.Printf("]")
		t.params[1].Emit(w)
	case "func":
		w.Printf("(")
		for i, p := range t.funcParams {
			p.Emit(w)
			if i < len(t.funcParams)-1 {
				w.Printf(", ")
			}
		}
		w.Printf(") (")
		for i, p := range t.funcRetParams {
			p.Emit(w)
			if i < len(t.funcRetParams)-1 {
				w.Printf(", ")
			}
		}
		w.Printf(")")

	default:
		for _, p := range t.params {
			p.Emit(w)
			w.Printf(" ")
		}
	}

}
