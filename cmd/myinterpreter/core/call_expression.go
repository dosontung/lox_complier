package core

import "github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/tokenize"

type CallExpression struct {
	Callee Expression
	Paren  *tokenize.Token
	Params []Expression
}

var _ Expression = (*CallExpression)(nil)

func (expr *CallExpression) Type() ExpressionType { return GROUPING }

func (expr *CallExpression) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitCallExpr(expr)
}
