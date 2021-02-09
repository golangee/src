package ast

import "go/token"

// TokenKind determines which kind of token, like INT, FLOAT, IMAG, CHAR, or STRING is meant.
type TokenKind int

const (
	TokenInt    TokenKind = TokenKind(token.INT)
	TokenFloat  TokenKind = TokenKind(token.FLOAT)
	TokenImag   TokenKind = TokenKind(token.IMAG)
	TokenChar   TokenKind = TokenKind(token.CHAR)
	TokenString TokenKind = TokenKind(token.STRING)
)

// A BasicLit represents a literal of a basic type.
type BasicLit struct {
	Kind  TokenKind
	Value string // the actual literal string, strings and chars must include the according escapes.
	Obj
}
