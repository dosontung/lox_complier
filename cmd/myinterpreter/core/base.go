package core

type ExpressionType string
type StatementType string

const (
	LITERAL  ExpressionType = "LITERAL"
	UNARY    ExpressionType = "UNARY"
	BINARY   ExpressionType = "BINARY"
	GROUPING ExpressionType = "GROUPING"
)

const (
	PRINT StatementType = "LITERAL"
)

type Expression interface {
	Type() ExpressionType
	Accept(visitor ExprVisitor) interface{}
}

type ExprVisitor interface {
	VisitBinaryExpr(*BinaryExpression) interface{}
	VisitGroupingExpr(*GroupExpression) interface{}
	VisitLiteralExpr(*LiteralExpression) interface{}
	VisitUnaryExpr(*UnaryExpression) interface{}
}

type Statement interface {
	Type() StatementType
	Accept(visitor StatementVisitor) interface{}
}

type StatementVisitor interface {
	VisitPrintStmt(statement *PrintStatement) interface{}
}
