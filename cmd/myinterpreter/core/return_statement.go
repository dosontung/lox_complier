package core

type ReturnStatement struct {
	Expr Expression
}

func (p *ReturnStatement) Accept(visitor StatementVisitor) interface{} {
	visitor.VisitReturnStmt(p)
	return nil
}

func (p *ReturnStatement) Type() StatementType {
	return RETURN
}

var _ Statement = (*ReturnStatement)(nil)
