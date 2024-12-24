package core

import (
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/tokenize"
)

type FuncStatement struct {
	Name     *tokenize.Token
	Params   []*tokenize.Token
	Body     []Statement
	LocalVar interface{} // only use for local func
}

func (p *FuncStatement) Accept(visitor StatementVisitor) interface{} {
	visitor.VisitFuncStmt(p)
	return nil
}

func (p *FuncStatement) Type() StatementType {
	return FUNC
}

func (p *FuncStatement) SetLocalVar(Env interface{}) {
	p.LocalVar = Env
}

func (p *FuncStatement) Clone() *FuncStatement {

	return &FuncStatement{
		Name:   p.Name,
		Params: p.Params,
		Body:   p.Body,
	}
}

var _ Statement = (*FuncStatement)(nil)
