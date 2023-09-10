package parser

import (
	"github.com/fabiante/monkeylang/ast"
	"github.com/fabiante/monkeylang/lexer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParser_ParseProgram(t *testing.T) {
	t.Run("returns error in invalid let statement", func(t *testing.T) {
		input := `let x 5;`

		lex := lexer.NewLexer(input)
		par := NewParser(lex)

		_ = par.ParseProgram()
		require.Len(t, par.Errors(), 1, "unexpected error count")
	})

	t.Run("let statement", func(t *testing.T) {
		input := `let x = 5;let y= 10;let foobar = 123;`

		tests := []struct {
			expectedIdentifier string
		}{
			{"x"},
			{"y"},
			{"foobar"},
		}

		lex := lexer.NewLexer(input)
		par := NewParser(lex)

		program := par.ParseProgram()
		requireNoParserErrors(t, par)
		require.NotNil(t, program)
		require.Len(t, program.Statements, len(tests))

		for i, test := range tests {
			stmt := program.Statements[i]
			assertLetStatement(t, stmt, test.expectedIdentifier)
		}
	})

	t.Run("return statement", func(t *testing.T) {
		input := `return 12 + 5;`

		lex := lexer.NewLexer(input)
		par := NewParser(lex)

		program := par.ParseProgram()
		requireNoParserErrors(t, par)
		require.NotNil(t, program)
		require.Len(t, program.Statements, 1)

		stmt := program.Statements[0]
		assertReturnStatement(t, stmt)
	})

	t.Run("identifier expression", func(t *testing.T) {
		input := `foobar;`

		lex := lexer.NewLexer(input)
		par := NewParser(lex)

		program := par.ParseProgram()
		requireNoParserErrors(t, par)
		require.NotNil(t, program)
		require.Len(t, program.Statements, 1)

		stmt := program.Statements[0]
		stmtExpression, ok := stmt.(*ast.ExpressionStatement)
		require.True(t, ok, "stmt has unexpected type %T", stmt)

		identifier, ok := stmtExpression.Expression.(*ast.Identifier)
		require.True(t, ok, "stmt expression has unexpected type %T", stmt)

		assert.Equal(t, "foobar", identifier.Value)
		assert.Equal(t, "foobar", identifier.TokenLiteral())
	})
}

func assertLetStatement(t *testing.T, node ast.Statement, name string) {
	require.NotNil(t, node)
	assert.Equal(t, "let", node.TokenLiteral())

	stmt, ok := node.(*ast.LetStatement)
	require.True(t, ok, "node is not of expected type, got %T", node)

	assert.Equal(t, name, stmt.Name.Value)
	assert.Equal(t, name, stmt.Name.TokenLiteral())
}

func assertReturnStatement(t *testing.T, node ast.Statement) {
	require.NotNil(t, node)
	assert.Equal(t, "return", node.TokenLiteral())

	_, ok := node.(*ast.ReturnStatement)
	require.True(t, ok, "node is not of expected type, got %T", node)
}

func requireNoParserErrors(t *testing.T, p *Parser) {
	errs := p.Errors()

	if len(errs) == 0 {
		return
	}

	t.Errorf("parser had %d errors", len(errs))
	for i, err := range errs {
		t.Errorf("err %d: %s", i, err)
	}

	t.FailNow()
}
