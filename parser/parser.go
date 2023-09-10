package parser

import (
	"fmt"
	"github.com/fabiante/monkeylang/ast"
	"github.com/fabiante/monkeylang/lexer"
	"github.com/fabiante/monkeylang/token"
)

type Parser struct {
	lexer *lexer.Lexer

	currToken token.Token
	peekToken token.Token

	errors []string
}

func NewParser(lexer *lexer.Lexer) *Parser {
	p := &Parser{lexer: lexer, errors: make([]string, 0)}
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
		return nil // TODO: error?
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

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	peek := p.lexer.NextToken()
	p.peekToken = peek
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

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) currTokenIs(t token.TokenType) bool {
	return p.currToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}
