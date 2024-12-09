package core

import "github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/tokenize"

type VarDeclarationStatement struct {
	Name *tokenize.Token
	Expr Expression
}

func (p *VarDeclarationStatement) Accept(visitor StatementVisitor) interface{} {
	visitor.VisitVarDeclarationStmt(p)
	return nil
}

func (p *VarDeclarationStatement) Type() StatementType {
	return EXPRESSION
}

var _ Statement = (*VarDeclarationStatement)(nil)
