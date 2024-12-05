package core

type GroupExpression struct {
	Expr Expression
}

var _ Expression = (*GroupExpression)(nil)

func (expr *GroupExpression) Type() ExpressionType { return GROUPING }

func (expr *GroupExpression) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitGroupingExpr(expr)
}
