package ast

import "github.com/fabiante/monkeylang/token"

// Identifier is an Expression that represents an identifier.
//
// It is not a statement on purpose because in code like the following,
// identifiers produce a value:
//
//	let x = y // here y is an identifier which should produce a value
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) String() string {
	return i.Value
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) expressionNode() {}
