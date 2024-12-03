package parser

import "github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/tokenize"

type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator tokenize.Token
}

func (expr *BinaryExpression) Type() ExpressionType { return BINARY }

func (expr *BinaryExpression) Accept(visitor Visitor) interface{} {
	return visitor.VisitBinaryExpr(expr)
}
