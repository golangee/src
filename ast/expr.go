package ast

// An Expr is a specialised node. Due to our macros, this is expressiveness is useless and due to different
// languages wrong anyway, like assignments (which are statements in Go and expressions in Java).
type Expr interface {
	Node
	exprNode() // marker interface method
}

// Exprs is a builder helper for vargs.
func Exprs(expr ...Expr) []Expr {
	return expr
}
