package parser

import (
	"fmt"
	"math"
	"strings"
)

type ExpressionType string

const (
	LITERAL  ExpressionType = "LITERAL"
	UNARY    ExpressionType = "UNARY"
	BINARY   ExpressionType = "BINARY"
	GROUPING ExpressionType = "GROUPING"
)

type Visitor interface {
	VisitBinaryExpr(*BinaryExpression) interface{}
	VisitGroupingExpr(*GroupExpression) interface{}
	VisitLiteralExpr(*LiteralExpression) interface{}
	VisitUnaryExpr(*UnaryExpression) interface{}
}

type Expression interface {
	Type() ExpressionType
	Accept(visitor Visitor) interface{}
}

type VisitorImpl struct{}

func (v *VisitorImpl) parenthesize(name string, exprs ...Expression) string {
	var sb strings.Builder

	sb.WriteString("(")
	sb.WriteString(name)
	for _, expr := range exprs {
		sb.WriteString(" ")
		sb.WriteString(fmt.Sprintf("%v", expr.Accept(v))) // Pass the visitor if needed
	}
	sb.WriteString(")")

	return sb.String()
}

func (v *VisitorImpl) VisitBinaryExpr(expr *BinaryExpression) interface{} {
	return v.parenthesize(expr.Operator.Lexeme,
		expr.Left, expr.Right)
}

func (v *VisitorImpl) VisitGroupingExpr(expr *GroupExpression) interface{} {

	return v.parenthesize("group", expr.Expr)
}

func (v *VisitorImpl) VisitLiteralExpr(expr *LiteralExpression) interface{} {
	if number, ok := expr.Value.(float64); ok {
		if number == math.Trunc(number) {
			return fmt.Sprintf("%v.0", number)
		}
		return fmt.Sprintf("%v", number)
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (v *VisitorImpl) VisitUnaryExpr(expr *UnaryExpression) interface{} {
	return v.parenthesize(expr.Operator.Lexeme, expr.Right)
}
