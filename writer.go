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
	"go/format"
	"strings"
)

type Writer interface {
	Printf(str string, args ...interface{})
}

type BufferedWriter strings.Builder

func (b *BufferedWriter) Printf(format string, args ...interface{}) {
	if len(args) == 0 {
		(*strings.Builder)(b).WriteString(format)
		return
	}

	(*strings.Builder)(b).WriteString(fmt.Sprintf(format, args...))
}

func (b *BufferedWriter) String() string {
	return (*strings.Builder)(b).String()
}

func (b *BufferedWriter) Format() (string, error) {
	buf, err := format.Source([]byte(b.String()))
	if err != nil {
		return b.WithLinesNumbers(), err
	}

	return string(buf), nil
}

func (b *BufferedWriter) WithLinesNumbers() string {
	sb := &strings.Builder{}
	for i, line := range strings.Split(b.String(), "\n") {
		sb.WriteString(fmt.Sprintf("%4d: %s\n", i+1, line))
	}
	return sb.String()
}
