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

import (
	"fmt"
	"testing"
)

func TestDSL(t *testing.T) {
	fmt.Println(
		NewFile("testpkg").
			AddTypes(NewStruct("Abc").
				AddFields(
					NewField("Hello", NewTypeDecl("int")).
						AddTag("json", "hello", "omitempty").
						AddTag("xml", "xml_hello"),
				),
				NewInterface("Yoh").
					AddMethods(
						NewFunc("hello").AddParams(NewParameter("OhMy", NewTypeDecl("error"))).SetVariadic(true),
					),

				NewIntEnum("Status", "unknown", "running", "stopped").SetDoc("...is an enum test."),
				NewStringEnum("Status2", "unknown", "running", "stopped").SetDoc("...is an enum test."),
			).
			AddFuncs(
				NewFunc("MyPkgFunc"),
			).
			String(),
	)

	myIface := NewInterface("BookRepo").
		AddMethods(
			NewFunc("ReadAll").
				SetDoc("... tries to read all entities.").
				AddParams(
					NewParameter("offset", NewTypeDecl("int")),
					NewParameter("limit", NewTypeDecl("int64")),
				).
				AddResults(
					NewParameter("", NewSliceDecl(NewTypeDecl("Book"))),
					NewParameter("", NewTypeDecl("error")),
				),
			NewFunc("SideEffectOnly"),
		)

	fmt.Println(
		NewFile("testiface").
			AddTypes(
				Implement(myIface, true),
				ImplementMock(myIface),
			),
	)

	fmt.Println(NewSliceDecl(NewTypeDecl("net/http.Response")))

	fmt.Println(
		NewFile("testiface").AddTypes(
			NewInterface("ReaderWriter").
				AddEmbedded(
					NewTypeDecl("io.Writer"),
					NewTypeDecl("io.Reader"),
					NewTypeDecl("loan/core.Service"),
					NewTypeDecl("search/core.Service"),
				),

			NewStruct("Opts").AddFields(
				NewField("Identifier", NewTypeDecl("string")),
			).
				AddMethodFromJson("FromJson").
				AddMethodToJson("ToJson", true, false, true).
				AddMethodToJson("ToJson2", false, true, false),

		).AddVars(
			NewVar("x").SetRHS(NewBlock("5")),
			NewVar("y").SetRHS(NewBlock("7")),
			NewVar("z").SetRHS(NewBlock("\"hello\"")).SetDoc("...is the last variable."),
		).AddSideEffectImports("github.com/go-sql-driver/mysql"),


	)
}
