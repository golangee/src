package ast2

// A TypeParam represents a type parameter, generic or template:
//  * Go: only the following build-in types can have type parameters: map[X]Y, []T, [n]T, chan T
//  * Go draft: if this ever comes, we are prepared. The draft currently looks like SyncMap[X, Y]
//  * Java: just like List<T> or Map<X extends A, Y super B>)
type TypeParam struct {
	Identifier string // Identifier of the template name
	Obj
}
