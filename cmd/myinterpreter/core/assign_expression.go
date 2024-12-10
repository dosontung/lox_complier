package core

import "github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/tokenize"

type AssignExpression struct {
	Name *tokenize.Token
	Expr Expression
}

var _ Expression = (*AssignExpression)(nil)

func (expr *AssignExpression) Type() ExpressionType { return ASSIGN }

func (expr *AssignExpression) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitAssignExpr(expr)
}
