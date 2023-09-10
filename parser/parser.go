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
	panic("not implemented")
}

func (p *Parser) nextToken() {
	p.token = p.peekToken
	peek := p.lexer.NextToken()
	p.peekToken = peek
}
