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

type packageQualifier struct {
	Path string
	Name string
}

type importStatement struct {
	pkg  packageQualifier
	name string
}

func (stmt *importStatement) emit(b strings.Builder) {
	if stmt.name == "" {
		b.WriteString(stmt.pkg.Name)
	} else {
		b.WriteString(stmt.name)
	}

	b.WriteString(" \"")
	b.WriteString(stmt.pkg.Path)
	b.WriteString("\"\n")
}
