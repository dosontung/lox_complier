package evaluate

import (
	"fmt"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/parser"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/tokenize"
	"math"
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
	case tokenize.MINUS:
		return leftVal.(float64) - rightVal.(float64)
	case tokenize.PLUS:
		if _, ok := rightVal.(string); ok {
			return fmt.Sprintf("%s%s", leftVal.(string), rightVal.(string))
		}
		return leftVal.(float64) + rightVal.(float64)
	default: // tokenize.MINUS
		return nil
	}

}

func (v *Evaluator) VisitGroupingExpr(expr *parser.GroupExpression) interface{} {
	return expr.Expr.Accept(v)
}

func (v *Evaluator) VisitLiteralExpr(expr *parser.LiteralExpression) interface{} {
	strVal := expr.Value
	switch strVal.(type) {
	case string:
		return strVal
	case float64:
		if strVal == math.Trunc(strVal.(float64)) {
			return math.Trunc(strVal.(float64))
		}
		return strVal
	}
	return strVal
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
		return -expr.Right.Accept(v).(float64)
	}
}
