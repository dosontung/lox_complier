package core

import "github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/tokenize"

type VarExpression struct {
	Name *tokenize.Token
}

var _ Expression = (*VarExpression)(nil)

func (expr *VarExpression) Type() ExpressionType { return VARIABLE }

func (expr *VarExpression) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitVarExpr(expr)
}
