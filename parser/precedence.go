package parser

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
