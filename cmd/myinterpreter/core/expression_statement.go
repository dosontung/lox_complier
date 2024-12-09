package core

type ExpressionStatement struct {
	Expr Expression
}

func (p *ExpressionStatement) Accept(visitor StatementVisitor) interface{} {
	visitor.VisitExpressionStmt(p)
	return nil
}

func (p *ExpressionStatement) Type() StatementType {
	return EXPRESSION
}

var _ Statement = (*ExpressionStatement)(nil)
