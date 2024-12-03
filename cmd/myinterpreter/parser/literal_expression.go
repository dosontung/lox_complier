package parser

type LiteralExpression struct {
	Value interface{}
}

func (expr *LiteralExpression) Type() ExpressionType { return LITERAL }

func (expr *LiteralExpression) Accept(visitor Visitor) interface{} {
	return visitor.VisitLiteralExpr(expr)
}
