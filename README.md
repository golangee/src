# src
src is a specialized go code generator to support a few cases elegantly.
The API is intentionally not created to reflect any facette or possibility
of the language. However, it tries to support selected use cases way more
elegantly than it is possible with other existing solutions.

## alternatives
* Daves [jennifer](https://github.com/dave/jennifer)
* go [text template](https://golang.org/pkg/text/template/)
* a big [list](https://github.com/golang/go/wiki/GoGenerateTools) of many *go generate* tools

## example

An innocent code like the following
```go
package src_git

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
			).String(),
	)
}

```

results in the following full blown go code:
```go
package testpkg

import (
	strconv "strconv"
)

type Abc struct {
	Hello int `json:"hello,omitempty" xml:"xml_hello" `
}

type Yoh interface {
	hello(OhMy ...error)
}

// Status is an enum test.
type Status int16

// String returns the enums name.
func (e Status) String() string {
	switch e {
	case 1:
		return "unknown"
	case 2:
		return "running"
	case 3:
		return "stopped"
	default:
		return strconv.Itoa(int(e))
	}
}

// IsValid returns true, if the value represents a defined enum value.
func (e Status) IsValid() bool {
	return e >= 1 && e <= 3
}

const (
	// StatusUnknown represents the enum 'unknown'.
	StatusUnknown Status = 1
	// StatusRunning represents the enum 'running'.
	StatusRunning Status = 2
	// StatusStopped represents the enum 'stopped'.
	StatusStopped Status = 3
)

var (
	// StatusValues contains all valid enum states.
	StatusValues = []Status{0, 1, 2}
)

// Status2 is an enum test.
type Status2 string

// String returns the enums name.
func (e Status2) String() string {
	return string(e)
}

// IsValid returns true, if the value represents a defined enum value.
func (e Status2) IsValid() bool {
	switch e {
	case "unknown":
		return true
	case "running":
		return true
	case "stopped":
		return true
	default:
		return false
	}
}

const (
	// Status2Unknown represents the enum 'unknown'.
	Status2Unknown Status2 = "unknown"
	// Status2Running represents the enum 'running'.
	Status2Running Status2 = "running"
	// Status2Stopped represents the enum 'stopped'.
	Status2Stopped Status2 = "stopped"
)

var (
	// Status2Values contains all valid enum states.
	Status2Values = []Status2{"unknown", "running", "stopped"}
)

```