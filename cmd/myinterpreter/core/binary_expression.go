package core

import (
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/tokenize"
)

type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator *tokenize.Token
}

var _ Expression = (*BinaryExpression)(nil)

func (expr *BinaryExpression) Type() ExpressionType { return BINARY }

func (expr *BinaryExpression) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitBinaryExpr(expr)
}
