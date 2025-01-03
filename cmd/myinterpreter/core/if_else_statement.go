package core

type IFElseStatement struct {
	Expr       Expression
	ThenBranch Statement
	ElseBranch Statement
}

func (p *IFElseStatement) Accept(visitor StatementVisitor) interface{} {
	visitor.VisitIfElseStmt(p)
	return nil
}

func (p *IFElseStatement) Type() StatementType {
	return IF_ELSE
}

var _ Statement = (*IFElseStatement)(nil)
