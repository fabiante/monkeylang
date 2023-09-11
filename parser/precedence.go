package parser

import "github.com/fabiante/monkeylang/token"

type precedence int

const (
	_ precedence = iota
	lowest
	equals
	lessgreater
	sum
	product
	prefix
	call
)

// precedences maps token types to their precedence, used by Parser when parsing
// expressions.
var precedences = map[token.TokenType]precedence{
	token.EQ:       equals,
	token.NEQ:      equals,
	token.LT:       lessgreater,
	token.GT:       lessgreater,
	token.Plus:     sum,
	token.Minus:    sum,
	token.Slash:    product,
	token.Asterisk: product,
}
