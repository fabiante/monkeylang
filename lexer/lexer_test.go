package lexer

import (
	"github.com/fabiante/monkeylang/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLexer_NextToken(t *testing.T) {
	t.Run("special characters", func(t *testing.T) {
		input := `=+(){},;`

		tests := []struct {
			expectedType    token.TokenType
			expectedLiteral string
		}{
			{token.Assign, "="},
			{token.Plus, "+"},
			{token.LParen, "("},
			{token.RParen, ")"},
			{token.LBrace, "{"},
			{token.RBrace, "}"},
			{token.Comma, ","},
			{token.Semicolon, ";"},
			{token.EOF, ""},
		}

		lexer := NewLexer(input)

		for i, test := range tests {
			actual, err := lexer.NextToken()
			require.NoError(t, err, "error when parsing token %d", i)
			require.NotNil(t, actual, "parsing token %d returned nil", i)

			assert.Equal(t, test.expectedLiteral, actual.Literal, "unexpected token literal %d", i)
			assert.Equal(t, test.expectedType, actual.Type, "unexpected token type %d", i)
		}
	})
}
