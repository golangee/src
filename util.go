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
	"strings"
	"unicode"
)

// snakeCaseToCamelCase converts strings like "my_snake-case" into "MySnakeCase".
func snakeCaseToCamelCase(str string) string {
	sb := &strings.Builder{}
	nextUp := true
	for _, r := range str {
		if r == '-' || r == '_' {
			nextUp = true
			continue
		}

		if nextUp {
			sb.WriteRune(unicode.ToUpper(r))
			nextUp = false
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
