package ast2

import (
	"reflect"
)

var nodeType = reflect.TypeOf((*error)(nil)).Elem()

// ParentAs starts at the given node and walks up the parent hierarchy until the first found node is assignable to
// target or no more parents exists. Example:
//   mod := &yast.Mod{}
//   if ok := yast.ParentAs(&mod, someNode); ok{
//   ...
//   }
func ParentAs(node Node, target interface{}) bool {
	if target == nil {
		panic("ast: target cannot be nil")
	}

	val := reflect.ValueOf(target)
	typ := val.Type()
	if typ.Kind() != reflect.Ptr || val.IsNil() {
		panic("ast: target must be a non-nil pointer")
	}

	if e := typ.Elem(); e.Kind() != reflect.Interface && !e.Implements(nodeType) {
		panic("ast: *target must be interface or implement yast.Node")
	}

	targetType := typ.Elem()
	for node != nil {
		if reflect.TypeOf(node).AssignableTo(targetType) {
			val.Elem().Set(reflect.ValueOf(node))
			return true
		}

		node = node.Parent()
	}
	return false
}

func assertNotAttached(n Node) {
	if n.Parent() != nil {
		panic("assert: node " + reflect.TypeOf(n).String() + " is already attached to " + reflect.TypeOf(n.Parent()).String())
	}
}

func assertSettableParent(node Node) SettableParent {
	if sp, ok := node.(SettableParent); ok {
		return sp
	} else {
		panic("assert: node must be a SettableParent: " + reflect.TypeOf(node).String())
	}
}
