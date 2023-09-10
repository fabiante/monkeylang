package ast

import (
	"github.com/fabiante/monkeylang/token"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProgram_String(t *testing.T) {
	prog := NewProgram()
	prog.Statements = append(prog.Statements, &LetStatement{
		Token: token.Token{
			Type:    token.Let,
			Literal: "let",
		},
		Name: &Identifier{
			Token: token.Token{
				Type:    token.Identifier,
				Literal: "myVar",
			},
			Value: "myVar",
		},
		Value: &Identifier{
			Token: token.Token{
				Type:    token.Identifier,
				Literal: "anotherVar",
			},
			Value: "anotherVar",
		},
	})

	str := prog.String()

	require.Equal(t, "let myVar = anotherVar;", str)
}
