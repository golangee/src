package src_test

import (
	"fmt"
	"github.com/golangee/src/v2"
	"github.com/golangee/src/v2/golang"
	"github.com/golangee/src/v2/java"
	"github.com/golangee/src/v2/stdlib"
	"strconv"
	"testing"
)

func createFieldTable() []*src.Field {
	var res []*src.Field
	for _, t := range stdlib.Types {
		switch t {
		case stdlib.List:
			continue
		case stdlib.Map:
			continue
		case stdlib.Void:
			continue
		}

		visibility := src.Public
		res = append(res,
			src.NewField("field"+t[:len(t)-1], src.NewSimpleTypeDecl(src.Name(t))).
				SetDoc("...is a "+visibility.String()+" "+t).
				SetVisibility(visibility),
		)

		visibility = src.Private
		res = append(res,
			src.NewField("field"+t[:len(t)-1]+"Slice", src.NewSliceTypeDecl(src.NewSimpleTypeDecl(src.Name(t)))).
				SetDoc("...is a "+visibility.String()+" slice of "+t).
				SetVisibility(visibility),
		)

		visibility = src.PackagePrivate
		res = append(res,
			src.NewField("field"+t[:len(t)-1]+"Ptr", src.NewTypeDeclPtr(src.NewSimpleTypeDecl(src.Name(t)))).
				SetDoc("...is a "+visibility.String()+" pointer to "+t).
				SetVisibility(visibility),
		)
	}

	res = append(res,
		src.NewField("fieldMap", src.NewMapDecl(src.NewSimpleTypeDecl(stdlib.String), src.NewSimpleTypeDecl(stdlib.Int))).
			SetDoc("...is a string/int map").
			SetVisibility(src.Protected),
	)

	res = append(res,
		src.NewField("fieldList", src.NewListDecl(src.NewSimpleTypeDecl(stdlib.String))).
			SetDoc("...is a List of strings").
			SetVisibility(src.Public),
	)

	res = append(res,
		src.NewField("fieldChan", src.NewChanTypeDecl(src.NewSimpleTypeDecl(stdlib.String))).
			SetDoc("...is a channel of strings").
			SetVisibility(src.Public),
	)

	res = append(res,
		src.NewField("fieldArray", src.NewArrayTypeDecl(5, src.NewSimpleTypeDecl(stdlib.Int))).
			SetDoc("...is an array with exact 5 int elements").
			AddAnnotations(
				src.NewAnnotation("javax.xml.ws.ServiceMode").SetDefault("javax.xml.ws.Service.Mode.PAYLOAD"),
			).
			SetVisibility(src.Public),
	)

	res = append(res,
		src.NewField("fieldComplex", src.NewTypeDeclPtr(src.NewMapDecl(src.NewTypeDeclPtr(src.NewSimpleTypeDecl(stdlib.String)), src.NewListDecl(src.NewListDecl(src.NewSimpleTypeDecl(stdlib.Int)))))).
			SetDoc("...is a pointer to a map with keys, which are pointers to a string and values are a list of list of integer").
			AddAnnotations(
				src.NewAnnotation("javax.xml.ws.WebServiceRef").
					SetValue("mappedName", strconv.Quote("abc mapped name")).
					SetValue("type", "String.class"),
			).
			SetVisibility(src.Public),
	)

	res = append(res,
		src.NewField("fieldFunc",
			src.NewFuncTypeDecl().
				AddInputParams(
					src.NewParam("limit", src.NewSimpleTypeDecl(stdlib.Int)),
					src.NewParam("offset", src.NewSimpleTypeDecl(stdlib.Int32)),
				).
				AddOutputParams(
					src.NewParam("res", src.NewSliceTypeDecl(src.NewSimpleTypeDecl(stdlib.String))),
					src.NewParam("err", src.NewSimpleTypeDecl(stdlib.Error)),
				)).
			SetDoc("...is a function pointer.").
			SetVisibility(src.Public),
	)

	return res
}

func NewTranspilerModel() *src.Module {
	mod := src.NewModule().AddPackages(
		src.NewPackage("github.com/golangee/src/example", "pexample").
			SetDocPreamble("Code generated by golangee/architecture. DO NOT EDIT.").
			SetDoc("... is a cool package showing the transpiler possibilities.\nAnother line of important text.").
			AddSrcFiles(
				src.NewSrcFile("example").
					AddTypes(
						src.NewStruct("Test").
							SetDoc("...is a simple example of defining a class or struct.\n\n    Can we have newlines?").
							SetFinal(true).
							SetStatic(true).
							AddFields(createFieldTable()...).
							AddMethods(
								src.NewFunc("HelloWorld").
									SetDoc("...says hello world").
									SetParams(src.NewParam("args", src.NewSimpleTypeDecl(stdlib.String))).
									SetVariadic(true),
							),

						src.NewInterface("MyInterface").
							SetDoc("...shows how to use interfaces.").
							AddAnnotations(
								src.NewAnnotation("javax.persistence.Entity").SetDefault(`"testEntity"`),
								src.NewAnnotation("javax.persistence.Table").SetDefault(`"testTable"`),
							).
							AddTypes(
								src.NewStruct("ANestedType").
									SetDoc("... a cooler nested type").
									AddFields(
										src.NewField("aName", src.NewSimpleTypeDecl(stdlib.Int)),
									),
							).
							AddMethods(
								src.NewFunc("sayHello").
									SetDoc("...says hello to the world.").
									SetVisibility(src.Private),
								src.NewFunc("sayHello2").
									SetDoc("...says hello to the world.").
									AddResults(
										src.NewParam("", src.NewSimpleTypeDecl(stdlib.Void)),
									),
								src.NewFunc("sayHello3").
									SetDoc("...says hello to the world.").
									AddAnnotations(
										src.NewAnnotation("javax.persistence.Entity").SetDefault(`"testEntity"`),
										src.NewAnnotation("javax.persistence.Table").SetDefault(`"testTable"`),
									).
									AddParams(
										src.NewParam("limit", src.NewSimpleTypeDecl(stdlib.Int)).
											AddAnnotations(src.NewAnnotation("javax.persistence.Entity").SetDefault(`"testEntity"`)).
											SetDoc("...provides the maximum amount of stuff returned."),
										src.NewParam("offset", src.NewSimpleTypeDecl(stdlib.Int64)).
											SetDoc("...declares the row from which to start the return."),
									).
									AddResults(
										src.NewParam("", src.NewSimpleTypeDecl(stdlib.String)).
											SetDoc("...is a weired string result."),
										src.NewParam("", src.NewSimpleTypeDecl("Exception")).
											SetDoc("...if something went seriously wrong."),
									),
							),

					).AddFunctions(
					src.NewFunc("aUtilityFunc").
						SetDoc("...is a stateless function.").
						SetBody(src.NewBlock().
							L(src.Name("fmt.Println")),
					),
				),
			),
	)
	return mod
}

func TestJava(t *testing.T) {
	file := NewTranspilerModel()
	renderedFiles, err := java.Render(file)
	if err != nil {
		fmt.Println(err)
		for _, f := range renderedFiles {
			if f.Error != nil {
				fmt.Println(f.Error)
				fmt.Println(string(f.Buf))
			}
		}
		t.Fatal()
	}

	for _, f := range renderedFiles {
		fmt.Println(f.AbsoluteFileName)
		fmt.Println(string(f.Buf))
	}

}

func TestGo(t *testing.T) {
	file := NewTranspilerModel()
	renderedFiles, err := golang.Render(file)
	if err != nil {
		fmt.Println(err)
		for _, f := range renderedFiles {
			if f.Error != nil {
				fmt.Println(f.Error)
				fmt.Println(string(f.Buf))
			}
		}
		t.Fatal()
	}

	for _, f := range renderedFiles {
		fmt.Println(f.AbsoluteFileName)
		fmt.Println(string(f.Buf))
	}

}

func TestTranspiler(t *testing.T) {
	TestGo(t)
	TestJava(t)
}
