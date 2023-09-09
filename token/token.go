package token

type TokenType int

const (
	Illegal TokenType = iota
	EOF
	// Identifier is a user-defined identifier. This is the opposite
	// from keywords of the language.
	Identifier
	Int
	Assign
	Plus
	Comma
	Semicolon
	LParen
	RParen
	LBrace
	RBrace
	Func
	Let
)

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"let": Let,
	"fn":  Func,
}

func LookupIdentifier(literal string) TokenType {
	if t, ok := keywords[literal]; ok {
		return t
	}
	return Identifier
}
