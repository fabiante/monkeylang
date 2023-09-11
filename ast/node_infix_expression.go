package ast

import (
	"bytes"
	"github.com/fabiante/monkeylang/token"
)

type InfixExpression struct {
	Token    token.Token
	Operator string
	Left     Expression
	Right    Expression
}

func (p *InfixExpression) TokenLiteral() string {
	return p.Token.Literal
}

func (p *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(p.Left.String())
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")
	return out.String()
}

func (p *InfixExpression) expressionNode() {}
