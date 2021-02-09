package golang

import (
	"fmt"
	. "github.com/golangee/src/ast"
	"testing"
)

func TestRenderer_Render(t *testing.T) {
	prj := newProject()
	renderer := NewRenderer(Options{})
	artifact, err := renderer.Render(prj)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(artifact)
}

func newProject() *Prj {
	return NewPrj("MyEpicProject").
		AddModules(
			NewMod("github.com/myproject/mymodule").
				SetLang(LangGo).
				AddPackages(
					NewPkg("github.com/myproject/mymodule/cmd/main").
						AddFiles(
							NewFile("main.go").
								AddTypes(
									NewStruct("HelloWorld"),
								),
						),
				),

		)
}
