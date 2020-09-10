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

type VarBuilder struct {
	doc      string
	name     string
	rhs      *Block
	typeDecl *TypeDecl
	parent   FileProvider
}

func NewVar(name string) *VarBuilder {
	return &VarBuilder{
		name: name,
	}
}

func (t *VarBuilder) SetRHS(block *Block) *VarBuilder {
	t.rhs = block
	t.rhs.onAttach(t)
	return t
}

func (t *VarBuilder) RHS() *Block {
	return t.rhs
}

func (t *VarBuilder) SetDoc(doc string) *VarBuilder {
	t.doc = doc
	return t
}

func (t *VarBuilder) Doc() string {
	return t.doc
}

func (t *VarBuilder) SetName(name string) *VarBuilder {
	t.name = name
	return t
}

func (t *VarBuilder) Name() string {
	return t.name
}

func (t *VarBuilder) SetType(decl *TypeDecl) *VarBuilder {
	t.typeDecl = decl
	return t
}

func (t *VarBuilder) Type() *TypeDecl {
	return t.typeDecl
}

func (t *VarBuilder) onAttach(parent FileProvider) {
	if t == nil {
		return
	}

	t.parent = parent
	if t.typeDecl != nil {
		t.typeDecl.onAttach(parent)
	}
}

func (t *VarBuilder) File() *FileBuilder {
	return t.parent.File()
}

func (t *VarBuilder) Emit(w Writer) {
	emitDoc(w, t.name, t.doc)
	if _, ok := t.parent.(*FileBuilder); ok {
		w.Printf("var ")
	}
	w.Printf(t.name)
	w.Printf(" ")
	if t.typeDecl != nil {
		t.typeDecl.Emit(w)
		w.Printf(" ")
	}
	w.Printf("= ")
	t.rhs.Emit(w)
}
