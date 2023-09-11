package ast

import "bytes"

type Program struct {
	Statements []Statement
}

func NewProgram() *Program {
	return &Program{
		Statements: make([]Statement, 0, 64),
	}
}

func (p *Program) statementNode() {}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, statement := range p.Statements {
		out.WriteString(statement.String())
	}

	return out.String()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
