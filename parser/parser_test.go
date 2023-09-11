package parser

import (
	"fmt"
	"github.com/fabiante/monkeylang/ast"
	"github.com/fabiante/monkeylang/lexer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strconv"
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
			assertLetStatement(t, test.expectedIdentifier, stmt)
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

	t.Run("integer literal expression", func(t *testing.T) {
		input := `5;`

		lex := lexer.NewLexer(input)
		par := NewParser(lex)

		program := par.ParseProgram()
		requireNoParserErrors(t, par)
		require.NotNil(t, program)
		require.Len(t, program.Statements, 1)

		stmt := program.Statements[0]
		stmtExpression, ok := stmt.(*ast.ExpressionStatement)
		require.True(t, ok, "stmt has unexpected type %T", stmt)

		assertIntegerLiteral(t, 5, stmtExpression.Expression)
	})

	t.Run("prefix operators", func(t *testing.T) {
		tests := []struct {
			input    string
			operator string
			value    int64
		}{
			{"-15;", "-", 15},
			{"!9;", "!", 9},
		}

		for _, test := range tests {
			t.Run(test.operator, func(t *testing.T) {
				lex := lexer.NewLexer(test.input)
				par := NewParser(lex)

				program := par.ParseProgram()
				requireNoParserErrors(t, par)
				require.NotNil(t, program)
				require.Len(t, program.Statements, 1)

				stmt := program.Statements[0]
				stmtExpression, ok := stmt.(*ast.ExpressionStatement)
				require.True(t, ok, "stmt has unexpected type %T", stmt)

				prefix, ok := stmtExpression.Expression.(*ast.PrefixExpression)
				require.True(t, ok, "prefix expression has unexpected type %T", stmt)

				assert.Equal(t, test.operator, prefix.Operator)
				assertIntegerLiteral(t, test.value, prefix.Right)
			})
		}
	})

	t.Run("infix operators", func(t *testing.T) {
		tests := []struct {
			input    string
			operator string
			left     int64
			right    int64
		}{
			{"5 + 6;", "+", 5, 6},
			{"5 - 6;", "-", 5, 6},
			{"5 * 6;", "*", 5, 6},
			{"5 / 6;", "/", 5, 6},
			{"5 > 6;", ">", 5, 6},
			{"5 < 6;", "<", 5, 6},
			{"5 == 6;", "==", 5, 6},
			{"5 != 6;", "!=", 5, 6},
		}

		for _, test := range tests {
			t.Run(test.operator, func(t *testing.T) {
				lex := lexer.NewLexer(test.input)
				par := NewParser(lex)

				program := par.ParseProgram()
				requireNoParserErrors(t, par)
				require.NotNil(t, program)
				require.Len(t, program.Statements, 1)

				stmt := program.Statements[0]
				stmtExpression, ok := stmt.(*ast.ExpressionStatement)
				require.True(t, ok, "stmt has unexpected type %T", stmt)

				prefix, ok := stmtExpression.Expression.(*ast.InfixExpression)
				require.True(t, ok, "prefix expression has unexpected type %T", stmt)

				assert.Equal(t, test.operator, prefix.Operator)
				assertIntegerLiteral(t, test.left, prefix.Left)
				assertIntegerLiteral(t, test.right, prefix.Right)
			})
		}
	})

	t.Run("operator precedence", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{
				"a + b * c",
				"(a + (b * c))",
			},
			{
				"a + b / c",
				"(a + (b / c))",
			},
			{
				"a + b * c + d / e - f",
				"(((a + (b * c)) + (d / e)) - f)",
			},
			{
				"5 > 4 == 3 < 4",
				"((5 > 4) == (3 < 4))",
			},
			{
				"5 < 4 != 3 > 4",
				"((5 < 4) != (3 > 4))",
			},
			{
				"3 + 4 * 5 == 3 * 1 + 4 * 5",
				"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
			},
			{
				"3 + 4 * 5 == 3 * 1 + 4 * 5",
				"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
			},
		}

		for i, test := range tests {
			t.Run(fmt.Sprintf("tests[%d]", i), func(t *testing.T) {
				lex := lexer.NewLexer(test.input)
				par := NewParser(lex)

				program := par.ParseProgram()
				requireNoParserErrors(t, par)
				require.NotNil(t, program)
				require.Len(t, program.Statements, 1)

				stmt := program.Statements[0]
				stmtExpression, ok := stmt.(*ast.ExpressionStatement)
				require.True(t, ok, "stmt has unexpected type %T", stmt)

				assert.Equal(t, test.expected, stmtExpression.String())
				assert.Equal(t, test.expected, program.String())
			})
		}
	})
}

func assertLetStatement(t *testing.T, name string, node ast.Statement) {
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

func assertIntegerLiteral(t *testing.T, expected int64, node ast.Expression) {
	identifier, ok := node.(*ast.IntegerLiteral)
	require.True(t, ok, "node has unexpected type %T", node)

	assert.Equal(t, expected, identifier.Value)
	assert.Equal(t, strconv.FormatInt(expected, 10), identifier.TokenLiteral())
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
