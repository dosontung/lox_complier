package core

type ExpressionType string
type StatementType string

/*
literal        → NUMBER | STRING | "true" | "false" | "nil" ;
grouping       → "(" expression ")" ;
unary          → ( "-" | "!" ) expression ;
binary         → expression operator expression ;
operator       → "==" | "!=" | "<" | "<=" | ">" | ">="
               | "+"  | "-"  | "*" | "/" ;
*/

const (
	LITERAL  ExpressionType = "LITERAL"
	LOGICAL  ExpressionType = "LOGICAL"
	UNARY    ExpressionType = "UNARY"
	BINARY   ExpressionType = "BINARY"
	GROUPING ExpressionType = "GROUPING"
	VARIABLE ExpressionType = "VARIABLE"
	ASSIGN   ExpressionType = "ASSIGN"
)

const (
	PRINT           StatementType = "PRINT"
	EXPRESSION      StatementType = "EXPRESSION"
	VAR_DECLARATION StatementType = "VAR_DECLARATION"
	BLOCK           StatementType = "BLOCK"
	IF_ELSE         StatementType = "IF_ELSE"
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
	VisitVarExpr(*VarExpression) interface{}
	VisitAssignExpr(*AssignExpression) interface{}
	VisitLogicalExpr(*LogicalExpression) interface{}
}

type Statement interface {
	Type() StatementType
	Accept(visitor StatementVisitor) interface{}
}

type StatementVisitor interface {
	VisitPrintStmt(statement *PrintStatement) interface{}
	VisitExpressionStmt(statement *ExpressionStatement)
	VisitVarDeclarationStmt(statement *VarDeclarationStatement)
	VisitBlockStmt(statement *BlockStatement)
	VisitIfElseStmt(statement *IFElseStatement)
}
