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

type Parameter struct {
	name   string
	decl   *TypeDecl
	parent FileProvider
}

func NewParameter(name string, decl *TypeDecl) *Parameter {
	return &Parameter{
		name: name,
		decl: decl,
	}
}

func (p *Parameter) Clone() *Parameter {
	return NewParameter(p.name, p.decl.Clone())
}

func (p *Parameter) Name() string {
	return p.name
}

func (p *Parameter) Decl() *TypeDecl {
	return p.decl
}

func (p *Parameter) onAttach(parent FileProvider) {
	p.parent = parent
	p.decl.onAttach(parent)
}

func (p *Parameter) Emit(w Writer) {
	w.Printf(p.name)
	w.Printf(" ")
	p.decl.Emit(w)
}

func (p *Parameter) emitAsVariadic(w Writer) {
	w.Printf(p.name)
	w.Printf("...")
	p.decl.Emit(w)
}
