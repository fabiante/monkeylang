package lexer

import (
	"github.com/fabiante/monkeylang/token"
)

type Lexer struct {
	input string

	// pos is the current position in input (points to current char).
	pos int
	// nextPos is the current reading position in input (after current char).
	nextPos int

	// char is the current char under examination.
	//
	// Note: Since this is a byte, the lexer can only work with single-byte characters (ASCII).
	char byte
}

func NewLexer(input string) *Lexer {
	lexer := &Lexer{input: input}
	lexer.readChar() // advance to first char
	return lexer
}

func (l *Lexer) NextToken() (token.Token, error) {
	var t token.Token

	switch l.char {
	case '=':
		t = newToken(token.Assign, string(l.char))
	case '+':
		t = newToken(token.Plus, string(l.char))
	case '(':
		t = newToken(token.LParen, string(l.char))
	case ')':
		t = newToken(token.RParen, string(l.char))
	case '{':
		t = newToken(token.LBrace, string(l.char))
	case '}':
		t = newToken(token.RBrace, string(l.char))
	case ',':
		t = newToken(token.Comma, string(l.char))
	case ';':
		t = newToken(token.Semicolon, string(l.char))
	case 0:
		t = newToken(token.EOF, "")
	default:
		t = newToken(token.Illegal, string(l.char))
	}

	l.readChar()

	return t, nil
}

func newToken(tokenType token.TokenType, literal string) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: literal,
	}
}

func (l *Lexer) readChar() {
	if l.nextPos >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.nextPos]
	}

	l.pos = l.nextPos
	l.nextPos += 1
}
