package core

type LiteralExpression struct {
	Value interface{}
}

var _ Expression = (*LiteralExpression)(nil)

func (expr *LiteralExpression) Type() ExpressionType { return LITERAL }

func (expr *LiteralExpression) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLiteralExpr(expr)
}
