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
	l.skipWhitespace()

	var t token.Token

	t.Literal = string(l.char)

	switch l.char {
	case '=':
		t.Type = token.Assign
	case '+':
		t.Type = token.Plus
	case '(':
		t.Type = token.LParen
	case ')':
		t.Type = token.RParen
	case '{':
		t.Type = token.LBrace
	case '}':
		t.Type = token.RBrace
	case ',':
		t.Type = token.Comma
	case ';':
		t.Type = token.Semicolon
	case 0:
		t.Type = token.EOF
		t.Literal = ""
	default:
		if isLetter(l.char) {
			t.Literal = l.readIdentifier()
			t.Type = token.LookupIdentifier(t.Literal)
			return t, nil // readIdentifier already advanced chars
		} else {
			t = newToken(token.Illegal, string(l.char))
		}
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

func (l *Lexer) skipWhitespace() {
	c := l.char
	for c == ' ' || c == '\t' || c == '\n' || c == '\r' {
		l.readChar()
		c = l.char
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

func (l *Lexer) readIdentifier() string {
	pos := l.pos
	// read until a non-letter is encountered
	for isLetter(l.char) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func isLetter(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_'
}
