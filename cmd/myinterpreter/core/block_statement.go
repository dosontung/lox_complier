package core

type BlockStatement struct {
	Statements []Statement
}

func (p *BlockStatement) Accept(visitor StatementVisitor) interface{} {
	visitor.VisitBlockStmt(p)
	return nil
}

func (p *BlockStatement) Type() StatementType {
	return BLOCK
}

var _ Statement = (*BlockStatement)(nil)
