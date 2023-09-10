package parser

import (
	"github.com/fabiante/monkeylang/ast"
	"github.com/fabiante/monkeylang/lexer"
	"github.com/fabiante/monkeylang/token"
)

type Parser struct {
	lexer *lexer.Lexer

	token     token.Token
	peekToken token.Token
}

func NewParser(lexer *lexer.Lexer) *Parser {
	p := &Parser{lexer: lexer}
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
	switch p.token.Type {
	case token.Let:
		return p.parseLetStatement()
	default:
		return nil // TODO: error?
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{
		Token: p.token,
		Name:  nil,
		Value: nil,
	}

	if !p.expectPeek(token.Identifier) {
		return nil // TODO: error?
	}

	stmt.Name = &ast.Identifier{
		Token: p.token,
		Value: p.token.Literal,
	}

	if !p.expectPeek(token.Assign) {
		return nil // TODO: error?
	}

	// TODO: Parse expression instead of ignoring it and skipping to semicolon
	for !p.currTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) nextToken() {
	p.token = p.peekToken
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
		return false
	}
}

func (p *Parser) currTokenIs(t token.TokenType) bool {
	return p.token.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}
