package ast

import (
	"github.com/fabiante/monkeylang/token"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrefixExpression_String(t *testing.T) {
	ex := &PrefixExpression{
		Token:    token.Token{Type: token.Bang, Literal: "!"},
		Operator: "!",
		Right: &Identifier{
			Token: token.Token{Type: token.Identifier, Literal: "abc"},
			Value: "abc",
		},
	}

	assert.Equal(t, "(!abc)", ex.String())
}
