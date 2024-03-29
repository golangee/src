package golang_test

import (
	"fmt"
	. "github.com/golangee/src/ast"
	"github.com/golangee/src/golang"
	"github.com/golangee/src/stdlib"
	fmt2 "github.com/golangee/src/stdlib/fmt"
	"github.com/golangee/src/stdlib/lang"
	"github.com/golangee/src/stdlib/strings"
	"testing"
)

func TestRenderer_Render(t *testing.T) {
	prj := newProject()
	renderer := golang.NewRenderer(golang.Options{})
	artifact, err := renderer.Render(prj)
	if err != nil {
		fmt.Println(artifact)
		t.Fatal(err)
	}

	fmt.Println(artifact)
}

func testError() *File {
	errorFile := NewFile("errors.go")
	// some error macro stuff
	myErr := lang.NewError("Ticket").
		SetComment("...is the sum type of all domain errors.")

	notFound := lang.NewErrorCase("NotFound").
		SetComment("...describes that a domain entity has not been found where one has been expected.")
	myErr.AddCase(notFound)

	alreadyDeclared := lang.NewErrorCase("AlreadyDeclaredError").
		SetComment("...describes a situation where a domain entity has been found but that was unexpected.").
		AddProperty("id", NewSimpleTypeDecl(stdlib.UUID), "...is the affected id.").
		AddProperty("status", NewSimpleTypeDecl(stdlib.Int), "...is some secret status code.")
	myErr.AddCase(alreadyDeclared)

	errorFile.AddNodes(myErr.TypeDecl())

	errorFile.AddNodes(
		NewFunc("TestError").
			AddResults(
				NewParam("", NewSimpleTypeDecl("error")),
			).SetBody(
			NewBlock(
				alreadyDeclared.Check(lang.CheckExactBehavior, "err", "matchedErr", NewBlock()),
				alreadyDeclared.Check(lang.CheckSumBehavior, "err", "matchedErr2", NewBlock()),
				alreadyDeclared.Check(lang.CheckCaseBehavior, "err", "matchedErr3", NewBlock()),
				NewReturnStmt(notFound.Make()),
				NewReturnStmt(alreadyDeclared.Make(NewIdent("nil"), NewIdentLit("42"))),
			),
		),
	)

	return errorFile
}

func newProject() *Prj {
	preamble := "Code generated by golangee/architecture. DO NOT EDIT."

	return NewPrj("MyEpicProject").
		AddModules(
			NewMod("github.com/myproject/mymodule").
				SetLang(LangGo).
				SetOutputDirectory("my/cool/module").
				SetLangVersion(LangVersionGo16).
				Require("github.com/golangee/sql v0.0.0-20210531101020-33021aed64c2").
				AddPackages(
					NewPkg("github.com/myproject/mymodule").
						AddRawFiles(NewRawFile("test.txt", "text", []byte("are we root yet?"))),
					NewPkg("github.com/myproject/mymodule/cmd/myapp").
						SetName("main").
						SetPreamble(preamble).
						SetComment("...is the actual package doc.").
						AddFiles(
							NewFile("main.go").
								SetPreamble(preamble).
								SetComment("...is a funny package.").
								AddNodes(NewImport("_", "github.com/go-sql-driver/mysql").SetComment("imported for sql driver side effect")).
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
										AddEmbedded(NewSimpleTypeDecl("sync.Mutex")).
										AddMethods(
											NewFunc("SayHello").
												SetComment("...shouts it into the world.").
												SetBody(
													NewBlock(
														NewBlock().SetComment("this is a redundant block"),
													),
												),
											NewFunc("IFaceTryEvil").
												AddResults(
													NewParam("", NewSimpleTypeDecl("HelloIface")),
													NewParam("", NewSimpleTypeDecl(stdlib.Error)),
												).SetBody(NewBlock(
												lang.TryDefine(NewIdent("db"),
													lang.CallStatic("database/sql.Open", lang.CallIdent("opts", "DSN")),
													"asdf",
												),

												lang.TryDefine(nil,
													lang.CallStatic("database/sql.Open", lang.CallIdent("opts", "DSN")),
													"asdf",
												),

												lang.Sel("a", "b", "c"),
												lang.Term(),

												NewForStmt(
													NewSimpleAssign(NewIdent("i"), AssignDefine, NewIntLit(0)),
													NewBinaryExpr(NewIdent("i"), OpLess, NewIntLit(10)),
													NewUnaryExpr(NewIdent("i"), OpInc),
													NewBlock(),
												),
												lang.Term(),


												NewRangeStmt(
													NewIdent("idx"),
													NewIdent("val"),
													lang.CallStatic("strings.Split", NewStrLit("hello"), NewStrLit("l")),
													NewBlock(),
												),
												lang.Term(),

												NewRangeStmt(
													nil,
													NewIdent("val"),
													lang.CallStatic("strings.Split", NewStrLit("hello"), NewStrLit("l")),
													NewBlock(),
												),
												lang.Term(),

												NewRangeStmt(
													NewIdent("idx"),
													nil,
													lang.CallStatic("strings.Split", NewStrLit("hello"), NewStrLit("l")),
													NewBlock(),
												),
												lang.Term(),

												NewForStmt(nil, lang.CallIdent("rows", "Next"), nil, NewBlock()),

												NewTpl(`// having a hard days life
{
var (x string = "abc") // ugly
/* don't do this at home */
fmt.Println({{ .Use "unsafe.Pointer"}}(x))
fmt.Println({{.Get "var"}})
}
`).Put("var", "XYZ"),

											)),
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
													NewDeferStmt(lang.CallIdent("db", "Close")),
													lang.Term(),
													lang.Panic("should come here"),
												)),
										),
								).
								AddFuncs(
									NewFunc("globalFunc").
										SetComment("...is a package private function.").
										SetVisibility(PackagePrivate).
										AddErrorCaseRefs(
											lang.NewErrorCase("NotFound").SetComment("...is another not-foundable error type."),
										).
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

									NewVarDecl(
										NewSimpleAssign(NewIdent("v"), AssignSimple, NewBasicLit(TokenString, "`dude`")).
											SetComment("...is another var."),
									),

									NewVarDecl(
										NewSimpleAssign(NewIdent("v2"), AssignSimple, NewBasicLit(TokenString, "`dude2`")).
											SetComment("...is another var."),
										NewSimpleAssign(NewIdent("v3"), AssignSimple, NewBasicLit(TokenString, "`dude3`")).
											SetComment("...is another var."),
										NewParam("v4", NewSliceTypeDecl(NewSimpleTypeDecl(stdlib.String))),
									),


								),
							testError(),
						).AddRawFiles(
						NewRawTpl("makefile", "text/x-makefile", NewTpl(
							`lint:
	@command -v golangci-lint
.PHONY: lint

VERSION = '{{.Get "check"}}'
`,
						).Put("check", "1.2.3")),
					),
				),

		)
}
