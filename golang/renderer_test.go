package golang

import (
	"fmt"
	. "github.com/golangee/src/ast"
	"github.com/golangee/src/stdlib"
	fmt2 "github.com/golangee/src/stdlib/fmt"
	"github.com/golangee/src/stdlib/lang"
	"github.com/golangee/src/stdlib/strings"
	"testing"
)

func TestRenderer_Render(t *testing.T) {
	prj := newProject()
	renderer := NewRenderer(Options{})
	artifact, err := renderer.Render(prj)
	if err != nil {
		fmt.Println(artifact)
		t.Fatal(err)
	}

	fmt.Println(artifact)
}

func newProject() *Prj {
	preamble := "Code generated by golangee/architecture. DO NOT EDIT."

	return NewPrj("MyEpicProject").
		AddModules(
			NewMod("github.com/myproject/mymodule").
				SetLang(LangGo).
				AddPackages(
					NewPkg("github.com/myproject/mymodule/cmd/myapp").
						SetName("main").
						SetPreamble(preamble).
						SetComment("...is the actual package doc.").
						AddFiles(
							NewFile("main.go").
								SetPreamble(preamble).
								SetComment("...is a funny package.").
								AddTypes(
									NewInterface("HelloIface").
										SetComment("...says hello").
										AddMethods(
											NewFunc("Wayne").
												SetComment("...cares a lot.").
												AddParams(NewParam("hey", NewSimpleTypeDecl("string"))),
										),

									NewStruct("HelloWorld").
										SetComment("... shows a struct.").
										AddFields(
											NewField("Hello", NewSimpleTypeDecl(stdlib.String)).
												SetComment("...holds a hello string."),
											NewField("World", NewSimpleTypeDecl(stdlib.String)).
												SetComment("...holds a world string.").
												AddAnnotations(
													NewAnnotation("json").SetDefault("world"),
													NewAnnotation("db").SetDefault("hello_world"),
												),
										).
										AddMethods(
											NewFunc("SayHello").
												SetComment("...shouts it into the world.").
												SetBody(
													NewBlock(
														NewBlock().SetComment("this is a redundant block"),
													),
												),
											NewFunc("Hello2").
												SetComment("...is a more complex method.").
												AddParams(
													NewParam("hey", NewSimpleTypeDecl(stdlib.Int)).SetComment("...declares a number."),
													NewParam("ho", NewSimpleTypeDecl(stdlib.Float64)).SetComment("...declares a float."),
												).
												AddResults(
													NewParam("", NewSliceTypeDecl(NewSimpleTypeDecl(stdlib.String))).SetComment("...a list of strings."),
													NewParam("", NewSimpleTypeDecl(stdlib.String)).SetComment("...declares a number."),
													NewParam("", NewSimpleTypeDecl(stdlib.Error)).SetComment("...is returned if everything fails."),
												).
												SetRecName("a").
												SetBody(NewBlock(
													fmt2.Println(NewIdent("hey"), NewIdent("ho"), NewStrLit("hello world")), lang.Term(),
													lang.TryDefine(NewIdent("rows"), lang.CallStatic("sql.query"), "cannot query"),
													strings.NewStrBuilder("sb",
														NewStrLit("test"),
														NewSelExpr(NewIdent("h"), NewIdent("myAttr")),
														lang.Attr("myAttr"),
														lang.ToString(lang.Attr("otherAttr")),
														lang.Itoa(lang.Attr("someInt")),
													),
												)),
										),
								).
								AddFuncs(
									NewFunc("globalFunc").
										SetComment("...is a package private function.").
										SetVisibility(PackagePrivate).
										SetBody(NewBlock()),
								).
								AddNodes(
									NewConstDecl(
										NewSimpleAssign(NewIdent("X"), AssignSimple, NewBasicLit(TokenString, "`hello`")).
											SetComment("...is a constant."),
									),

									NewConstDecl(
										NewSimpleAssign(NewIdent("a"), AssignSimple, NewBasicLit(TokenString, "`world`")).
											SetComment("...is a another."),
										NewSimpleAssign(NewIdent("b"), AssignSimple, NewBasicLit(TokenIdent, "4")).
											SetComment("...another cool constant."),
									),
								),

						),
				),

		)
}
