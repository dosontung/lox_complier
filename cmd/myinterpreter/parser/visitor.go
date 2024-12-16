package parser

import (
	"fmt"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/core"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/interpreter"
	"math"
	"strings"
)

type VisitorImpl struct {
	env *interpreter.Environment
}

func (v *VisitorImpl) VisitLogicalExpr(expression *core.LogicalExpression) interface{} {
	//TODO implement me
	panic("implement me")
}

var _ core.ExprVisitor = &VisitorImpl{}

func (v *VisitorImpl) parenthesize(name string, exprs ...core.Expression) string {
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

func (v *VisitorImpl) VisitBinaryExpr(expr *core.BinaryExpression) interface{} {
	return v.parenthesize(expr.Operator.Lexeme,
		expr.Left, expr.Right)
}

func (v *VisitorImpl) VisitGroupingExpr(expr *core.GroupExpression) interface{} {

	return v.parenthesize("group", expr.Expr)
}

func (v *VisitorImpl) VisitLiteralExpr(expr *core.LiteralExpression) interface{} {
	if number, ok := expr.Value.(float64); ok {
		if number == math.Trunc(number) {
			return fmt.Sprintf("%v.0", number)
		}
		return fmt.Sprintf("%v", number)
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (v *VisitorImpl) VisitUnaryExpr(expr *core.UnaryExpression) interface{} {
	return v.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (v *VisitorImpl) VisitVarExpr(expr *core.VarExpression) interface{} {
	err, i := v.env.GetKey(expr.Name.Lexeme)
	if err != nil {
		return nil
	}
	return i
}

func (v *VisitorImpl) VisitAssignExpr(expr *core.AssignExpression) interface{} {
	// TODO: Implement me >>>
	return nil
}
