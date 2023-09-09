package token

type TokenType int

const (
	Illegal TokenType = iota
	EOF
	Identifier
	Int
	Assign
	Plus
	Comma
	Semicolon
	LParen
	RParen
	LBrace
	RBrace
	Func
	Let
)

type Token struct {
	Type    TokenType
	Literal string
}
