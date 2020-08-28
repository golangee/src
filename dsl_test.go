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
