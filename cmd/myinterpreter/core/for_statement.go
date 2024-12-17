package core

type ForStatement struct {
	VarStatment Statement
	Expr1       Expression
	Expr2       Expression
	Body        Statement
}

func (p *ForStatement) Accept(visitor StatementVisitor) interface{} {
	visitor.VisitForStmt(p)
	return nil
}

func (p *ForStatement) Type() StatementType {
	return FOR
}

var _ Statement = (*ForStatement)(nil)
