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

	t.Run("small program", func(t *testing.T) {
		input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
    x + y;
};

let result = add(five, ten);`

		tests := []struct {
			expectedType    token.TokenType
			expectedLiteral string
		}{
			{token.Let, "let"},
			{token.Identifier, "five"},
			{token.Assign, "="},
			{token.Int, "5"},
			{token.Semicolon, ";"},
			{token.Let, "let"},
			{token.Identifier, "ten"},
			{token.Assign, "="},
			{token.Int, "10"},
			{token.Semicolon, ";"},
			{token.Let, "let"},
			{token.Identifier, "add"},
			{token.Assign, "="},
			{token.Func, "fn"},
			{token.LParen, "("},
			{token.Identifier, "x"},
			{token.Comma, ","},
			{token.Identifier, "y"},
			{token.RParen, ")"},
			{token.LBrace, "{"},
			{token.Identifier, "x"},
			{token.Plus, "+"},
			{token.Identifier, "y"},
			{token.Semicolon, ";"},
			{token.RBrace, "}"},
			{token.Semicolon, ";"},
			{token.Let, "let"},
			{token.Identifier, "result"},
			{token.Assign, "="},
			{token.Identifier, "add"},
			{token.LParen, "("},
			{token.Identifier, "five"},
			{token.Comma, ","},
			{token.Identifier, "ten"},
			{token.RParen, ")"},
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

	t.Run("operations and comparators", func(t *testing.T) {
		input := `!-/*5;5 < 10 > 5`

		tests := []struct {
			expectedType    token.TokenType
			expectedLiteral string
		}{
			{token.Bang, "!"},
			{token.Minus, "-"},
			{token.Slash, "/"},
			{token.Asterisk, "*"},
			{token.Int, "5"},
			{token.Semicolon, ";"},
			{token.Int, "5"},
			{token.LT, "<"},
			{token.Int, "10"},
			{token.GT, ">"},
			{token.Int, "5"},
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

	t.Run("if else return", func(t *testing.T) {
		input := `if true {} else {return 5;}`

		tests := []struct {
			expectedType    token.TokenType
			expectedLiteral string
		}{
			{token.If, "if"},
			{token.True, "true"},
			{token.LBrace, "{"},
			{token.RBrace, "}"},
			{token.Else, "else"},
			{token.LBrace, "{"},
			{token.Return, "return"},
			{token.Int, "5"},
			{token.Semicolon, ";"},
			{token.RBrace, "}"},
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

	t.Run("bool comparisons", func(t *testing.T) {
		input := `a == a; a != a;`

		tests := []struct {
			expectedType    token.TokenType
			expectedLiteral string
		}{
			{token.Identifier, "a"},
			{token.EQ, "=="},
			{token.Identifier, "a"},
			{token.Semicolon, ";"},
			{token.Identifier, "a"},
			{token.NEQ, "!="},
			{token.Identifier, "a"},
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
