package parser

import "github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/tokenize"

type GroupExpression struct {
	Left     Expression
	Right    Expression
	Operator tokenize.Token
}

func (expr *GroupExpression) Type() ExpressionType { return GROUPING }

func (expr *GroupExpression) Accept(visitor Visitor) interface{} {
	return visitor.VisitGroupingExpr(expr)
}
