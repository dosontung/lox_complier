package evaluate

import (
	"fmt"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/parser"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/tokenize"
	"math"
	"strconv"
)

type Evaluator struct {
}

func (v *Evaluator) VisitBinaryExpr(expr *parser.BinaryExpression) interface{} {
	rightVal := expr.Right.Accept(v)
	leftVal := expr.Left.Accept(v)
	switch expr.Operator.Type {
	case tokenize.STAR:
		return rightVal.(float64) * leftVal.(float64)
	case tokenize.SLASH:
		return leftVal.(float64) / rightVal.(float64)
	default: // tokenize.MINUS
		return nil
	}

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
	return floatValue
}

func (v *Evaluator) VisitUnaryExpr(expr *parser.UnaryExpression) interface{} {
	strVal := expr.Right.Accept(v)
	switch expr.Operator.Type {
	case tokenize.BANG:
		if strVal == false || strVal == "false" || strVal == "nil" {
			return true
		} else {
			return false
		}
	default: // tokenize.MINUS
		return fmt.Sprintf("-%v", strVal)
	}
}
