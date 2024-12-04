package evaluate

import (
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/parser"
	"math"
	"strconv"
)

type Evaluator struct {
}

func (v *Evaluator) VisitBinaryExpr(expr *parser.BinaryExpression) interface{} {
	return nil
}

func (v *Evaluator) VisitGroupingExpr(expr *parser.GroupExpression) interface{} {
	return expr.Expr.Accept(v)
}

func (v *Evaluator) VisitLiteralExpr(expr *parser.LiteralExpression) interface{} {
	strVal := expr.Value.(string)
	floatValue, err := strconv.ParseFloat(strVal, 64)
	if err != nil {
		return strVal
	}
	if floatValue == math.Trunc(floatValue) {
		return math.Trunc(floatValue)
	}
	return strVal
}

func (v *Evaluator) VisitUnaryExpr(expr *parser.UnaryExpression) interface{} {
	return nil

}
