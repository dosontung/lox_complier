package core

type PrintStatement struct {
	Expr Expression
}

func (p *PrintStatement) Accept(visitor StatementVisitor) interface{} {
	return visitor.VisitPrintStmt(p)
}

func (p *PrintStatement) Type() StatementType {
	return PRINT
}

var _ Statement = (*PrintStatement)(nil)
