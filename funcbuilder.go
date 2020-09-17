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

import "strings"

type FuncBuilder struct {
	parent        FileProvider
	parentType    *TypeBuilder
	doc           string
	name          string
	recName       string
	isPtrReceiver bool
	params        []*Parameter
	results       []*Parameter
	body          []*Block
	variadic      bool
}

func NewFunc(name string) *FuncBuilder {
	b := &FuncBuilder{}
	b.SetName(name)
	return b
}

func (b *FuncBuilder) Doc() string {
	return b.doc
}

func (b *FuncBuilder) Name() string {
	return b.name
}

func (b *FuncBuilder) Params() []*Parameter {
	return b.params
}

func (b *FuncBuilder) CloneParams() []*Parameter {
	var r []*Parameter
	for _, p := range b.params {
		r = append(r, p.Clone())
	}

	return r
}

func (b *FuncBuilder) AddParams(params ...*Parameter) *FuncBuilder {
	b.params = append(b.params, params...)
	for _, p := range params {
		p.onAttach(b)
	}

	return b
}

func (b *FuncBuilder) AddResults(params ...*Parameter) *FuncBuilder {
	b.results = append(b.results, params...)
	for _, p := range params {
		p.onAttach(b)
	}

	return b
}

func (b *FuncBuilder) SetVariadic(v bool) *FuncBuilder {
	b.variadic = v
	return b
}

func (b *FuncBuilder) Variadic() bool {
	return b.variadic
}

func (b *FuncBuilder) Results() []*Parameter {
	return b.results
}

func (b *FuncBuilder) CloneResults() []*Parameter {
	var r []*Parameter
	for _, p := range b.results {
		r = append(r, p.Clone())
	}

	return r
}

func (b *FuncBuilder) AddBody(blocks ...*Block) *FuncBuilder {
	b.body = append(b.body, blocks...)
	for _, block := range blocks {
		block.onAttach(b)
	}
	return b
}

func (b *FuncBuilder) File() *FileBuilder {
	return b.parent.File()
}

func (b *FuncBuilder) onAttach(parent FileProvider) {
	b.parent = parent
	if parent, ok := parent.(*TypeBuilder); ok {
		b.parentType = parent

		if b.recName == "" && len(parent.name) > 0 {
			b.recName = strings.ToLower(parent.name[0:1])
		}
	}
}

func (b *FuncBuilder) SetName(name string) *FuncBuilder {
	b.name = name
	return b
}

func (b *FuncBuilder) SetReceiverName(name string) *FuncBuilder {
	b.recName = name
	return b
}

func (b *FuncBuilder) ReceiverName() string {
	return b.recName
}

func (b *FuncBuilder) SetPointerReceiver(isPtrRec bool) *FuncBuilder {
	b.isPtrReceiver = isPtrRec
	return b
}

func (b *FuncBuilder) SetDoc(doc string) *FuncBuilder {
	b.doc = doc
	return b
}

func (b *FuncBuilder) Emit(w Writer) {
	emitDoc(w, b.name, b.doc)

	if b.parentType != nil && b.parentType.iType == typeInterface {
		w.Printf("%s(", b.name)
		for i, p := range b.params {
			if i == len(b.params)-1 && b.variadic {
				p.emitAsVariadic(w)
			} else {
				p.Emit(w)
			}
			w.Printf(",")
		}
		w.Printf(")")

		w.Printf("(")
		for _, p := range b.results {
			p.Emit(w)
			w.Printf(",")
		}
		w.Printf(")")

		return
	}

	w.Printf("func ")
	if b.parentType != nil {
		ptrRec := ""
		if b.isPtrReceiver {
			ptrRec = "*"
		}
		w.Printf("(%s %s%s) ", b.recName, ptrRec, b.parentType.name)
	}

	w.Printf("%s(", b.name)
	for _, p := range b.params {
		p.Emit(w)
		w.Printf(",")
	}
	w.Printf(")")

	w.Printf("(")
	for _, p := range b.results {
		p.Emit(w)
		w.Printf(",")
	}
	w.Printf(")")

	w.Printf("{\n")
	for _, block := range b.body {
		block.Emit(w)
	}
	w.Printf("}\n")
}
