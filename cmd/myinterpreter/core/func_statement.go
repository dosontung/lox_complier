package core

import "github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/tokenize"

type FuncStatement struct {
	Name   *tokenize.Token
	Params []*tokenize.Token
	Body   []Statement
}

func (p *FuncStatement) Accept(visitor StatementVisitor) interface{} {
	visitor.VisitFuncStmt(p)
	return nil
}

func (p *FuncStatement) Type() StatementType {
	return FUNC
}

var _ Statement = (*FuncStatement)(nil)
