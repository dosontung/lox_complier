package core

import "github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/tokenize"

type LogicalExpression struct {
	Operator *tokenize.Token
	Left     Expression
	Right    Expression
}

var _ Expression = (*LogicalExpression)(nil)

func (expr *LogicalExpression) Type() ExpressionType { return LOGICAL }

func (expr *LogicalExpression) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLogicalExpr(expr)
}
