package parser

import (
	"fmt"
	"github.com/fabiante/monkeylang/ast"
	"github.com/fabiante/monkeylang/lexer"
	"github.com/fabiante/monkeylang/token"
	"strconv"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(left ast.Expression) ast.Expression
)

type Parser struct {
	lexer *lexer.Lexer

	currToken token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn

	errors []string
}

func NewParser(lexer *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:          lexer,
		errors:         make([]string, 0),
		prefixParseFns: make(map[token.TokenType]prefixParseFn),
		infixParseFns:  make(map[token.TokenType]infixParseFn),
	}

	p.registerPrefixParseFn(token.Identifier, p.parseIdentifier)
	p.registerPrefixParseFn(token.Int, p.parseIntLiteral)
	p.registerPrefixParseFn(token.Bang, p.parsePrefixExpression)
	p.registerPrefixParseFn(token.Minus, p.parsePrefixExpression)

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	prog := ast.NewProgram()

	for !p.currTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			prog.Statements = append(prog.Statements, stmt)
		}
		p.nextToken()
	}

	return prog
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.Let:
		return p.parseLetStatement()
	case token.Return:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{
		Token: p.currToken,
		Name:  nil,
		Value: nil,
	}

	if !p.expectPeek(token.Identifier) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}

	if !p.expectPeek(token.Assign) {
		return nil
	}

	// TODO: Parse expression instead of ignoring it and skipping to semicolon
	for !p.currTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{
		Token:       p.currToken,
		ReturnValue: nil,
	}

	// TODO: Parse expression instead of ignoring it and skipping to semicolon
	for !p.currTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{
		Token:      p.currToken,
		Expression: p.parseExpression(lowest),
	}

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

// parseExpression parses an expression starting with the currToken and the given precedence.
func (p *Parser) parseExpression(precedence precedence) ast.Expression {
	prefix := p.prefixParseFns[p.currToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.currToken.Type)
		return nil
	}

	leftExp := prefix()

	return leftExp
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	peek := p.lexer.NextToken()
	p.peekToken = peek
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}
}

func (p *Parser) parseIntLiteral() ast.Expression {
	literal := p.currToken.Literal

	value, err := strconv.ParseInt(literal, 0, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("could not parse %s as integer", literal))
		return nil
	}

	return &ast.IntegerLiteral{
		Token: p.currToken,
		Value: value,
	}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
		Right:    nil,
	}

	p.nextToken()

	exp.Right = p.parseExpression(prefix)

	return exp
}

// expectPeek checks if the next token is of the given type,
//
// If it is, it will advance to the next token by calling nextToken and return true.
//
// If it isn't, it will return false.
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t token.TokenType) {
	err := fmt.Sprintf("expected token type %d, got %d instead", t, p.peekToken.Type)
	p.errors = append(p.errors, err)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	err := fmt.Sprintf("no prefix parse fn for token type %d", t)
	p.errors = append(p.errors, err)
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) currTokenIs(t token.TokenType) bool {
	return p.currToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) registerPrefixParseFn(t token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[t] = fn
}

func (p *Parser) registerInfixParseFn(t token.TokenType, fn infixParseFn) {
	p.infixParseFns[t] = fn
}

func (p *Parser) peekPrecedence() precedence {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return lowest
}

func (p *Parser) currPrecedence() precedence {
	if p, ok := precedences[p.currToken.Type]; ok {
		return p
	}
	return lowest
}
