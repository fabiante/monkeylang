package ast

type Node interface {
	// TokenLiteral is used only for debugging / testing
	TokenLiteral() string
	String() string
}

// A Statement is something that does not produce a value.
type Statement interface {
	Node
	statementNode()
}

// An Expression is something that will produce a value.
type Expression interface {
	Node
	expressionNode()
}
