package core

type ExpressionStatement struct {
	Expr Expression
}

func (p *ExpressionStatement) Accept(visitor StatementVisitor) interface{} {
	return visitor.VisitExpressionStmt(p)
}

func (p *ExpressionStatement) Type() StatementType {
	return EXPRESSION
}

var _ Statement = (*ExpressionStatement)(nil)
