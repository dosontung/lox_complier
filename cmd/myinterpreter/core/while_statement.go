package core

type WhileStatement struct {
	Expr Expression
	Body Statement
}

func (p *WhileStatement) Accept(visitor StatementVisitor) interface{} {
	visitor.VisitWhileStmt(p)
	return nil
}

func (p *WhileStatement) Type() StatementType {
	return WHILE
}

var _ Statement = (*WhileStatement)(nil)
