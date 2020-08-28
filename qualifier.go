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


// A Qualifier is a <path>.<name> string, e.g. github.com/myproject/mymod/mypath.MyType. Universe types have a
// leading dot (or no dot at all), like .int or .float32 or .error.
// It does not carry any information about the actual package name, so it can only be used in an explicitly named import
// context, which is sufficient per definition. Generic types cannot be expressed and must use the TypeDecl struct.
type Qualifier string

func (q Qualifier) Name() string {
	i := strings.LastIndex(string(q), ".")
	if i == -1 {
		return string(q)
	}

	return string(q[i+1:])
}

func (q Qualifier) Path() string {
	i := strings.LastIndex(string(q), ".")
	if i == -1 {
		return ""
	}

	return string(q[:i])
}
