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

type Const struct {
	doc      string
	name     string
	typeDecl *TypeDecl //may be nil for untyped
	value    string
	parent   FileProvider
}

func NewConst(name, rawValue string) *Const {
	return &Const{
		name:  name,
		value: rawValue,
	}
}

func NewTypedConst(name string, decl *TypeDecl, value string) *Const {
	return &Const{
		name:     name,
		typeDecl: decl,
		value:    value,
	}
}

func (t *Const) SetDoc(doc string) *Const {
	t.doc = doc
	return t
}

func (t *Const) Doc() string {
	return t.doc
}

func (t *Const) SetName(name string) *Const {
	t.name = name
	return t
}

func (t *Const) Name() string {
	return t.name
}

func (t *Const) SetType(decl *TypeDecl) *Const {
	t.typeDecl = decl
	return t
}

func (t *Const) Type() *TypeDecl {
	return t.typeDecl
}

func (t *Const) onAttach(parent FileProvider) {
	if t == nil {
		return
	}

	t.parent = parent
	if t.typeDecl != nil {
		t.typeDecl.onAttach(t)
	}
}

func (t *Const) File() *FileBuilder {
	return t.parent.File()
}

func (t *Const) Emit(w Writer) {
	emitDoc(w, t.name, t.doc)
	if t.typeDecl == nil {
		w.Printf("%s = %s\n", t.name, t.value)
	} else {
		w.Printf("%s ", t.name)
		t.typeDecl.Emit(w)
		w.Printf(" = %s\n", t.value)
	}
}
