package parser

import "github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/tokenize"

type UnaryExpression struct {
	Operator *tokenize.Token
	Right    Expression
}

func (expr *UnaryExpression) Type() ExpressionType { return UNARY }

func (expr *UnaryExpression) Accept(visitor Visitor) interface{} {
	return visitor.VisitUnaryExpr(expr)
}
