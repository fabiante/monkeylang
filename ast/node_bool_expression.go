package ast

import "github.com/fabiante/monkeylang/token"

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (b *BooleanLiteral) TokenLiteral() string {
	return b.Token.Literal
}

func (b *BooleanLiteral) String() string {
	return b.Token.Literal
}

func (b *BooleanLiteral) expressionNode() {}
