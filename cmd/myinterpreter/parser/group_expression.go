package parser

type GroupExpression struct {
	Expr Expression
}

func (expr *GroupExpression) Type() ExpressionType { return GROUPING }

func (expr *GroupExpression) Accept(visitor Visitor) interface{} {
	return visitor.VisitGroupingExpr(expr)
}
