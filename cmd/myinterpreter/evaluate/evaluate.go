package evaluate

import "github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/parser"

type Evaluator struct {
}

func (v *Evaluator) VisitBinaryExpr(expr *parser.BinaryExpression) interface{} {
	return nil
}

func (v *Evaluator) VisitGroupingExpr(expr *parser.GroupExpression) interface{} {
	return nil
}

func (v *Evaluator) VisitLiteralExpr(expr *parser.LiteralExpression) interface{} {
	return expr.Value
}

func (v *Evaluator) VisitUnaryExpr(expr *parser.UnaryExpression) interface{} {
	return nil

}
