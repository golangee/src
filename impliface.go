package src

import "strconv"

// Implement takes the interface type and returns a struct type.
// Each method is imported but the body is empty, so that
// the result will probably not compile.
func Implement(iface *TypeBuilder, ptrReceiver bool) *TypeBuilder {
	strct := NewStruct(iface.name + "Impl").
		SetDoc("...is an implementation of " + iface.name + ".\n" + iface.doc)
	for _, fun := range iface.Methods() {
		strct.AddMethods(
			NewFunc(fun.name).
				SetDoc(fun.Doc()).
				AddParams(fun.Params()...).
				AddResults(fun.Results()...).
				SetPointerReceiver(ptrReceiver),
		)
	}

	return strct
}

// ImplementMock takes the interface type and returns a struct type.
// Each method is delegated to a public member func which can be set
// as desired.
func ImplementMock(iface *TypeBuilder) *TypeBuilder {
	strct := NewStruct(iface.name + "Mock").
		SetDoc("...is a mock implementation of " + iface.name + ".\n" + iface.doc)
	for _, f := range iface.Methods() {
		fieldName := f.name + "Func"
		strct.AddFields(
			NewField(fieldName, NewFuncDecl(f.Params(), f.Results())).
				SetDoc("... mocks the " + f.Name() + " function."),
		)
		fun := NewFunc(f.name).
			SetDoc(f.Doc()).
			SetPointerReceiver(false).
			SetReceiverName("m")

		callParamList := ""
		for i, p := range f.Params() {
			name := p.name
			if name == "" {
				name = "p" + strconv.Itoa(i)
			}
			fun.AddParams(NewParameter(name, p.Decl().Clone()))
			callParamList += name
			if i < len(f.params)-1 {
				callParamList += ", "
			}
		}

		for _, p := range f.Results() {
			fun.AddResults(NewParameter(p.Name(), p.Decl().Clone()))
		}

		fun.AddBody(NewBlock().
			If(fun.ReceiverName()+"."+fieldName+"!=nil",
				NewBlock().Add("return ", fun.ReceiverName(), ".", fun.Name(), "(", callParamList, ")"),
			).
			AddLine("panic(\"mock not available: " + fun.Name() + "\")"),
		)
		strct.AddMethods(fun)
	}

	return strct
}
