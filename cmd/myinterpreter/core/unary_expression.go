package core

import (
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/tokenize"
)

type UnaryExpression struct {
	Operator *tokenize.Token
	Right    Expression
}

var _ Expression = (*UnaryExpression)(nil)

func (expr *UnaryExpression) Type() ExpressionType { return UNARY }

func (expr *UnaryExpression) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitUnaryExpr(expr)
}
