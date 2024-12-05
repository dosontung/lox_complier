package evaluate

import (
	"fmt"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/errors"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/parser"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/tokenize"
	"math"
	"os"
)

type Evaluator struct {
}

func (v *Evaluator) VisitBinaryExpr(expr *parser.BinaryExpression) interface{} {
	rightVal := expr.Right.Accept(v)
	leftVal := expr.Left.Accept(v)

	switch expr.Operator.Type {
	case tokenize.STAR:
		if _, ok := leftVal.(string); ok {
			v.raiseError(errors.OperandMustBeNumber)
		}
		if _, ok := rightVal.(string); ok {
			v.raiseError(errors.OperandMustBeNumber)
		}
		return rightVal.(float64) * leftVal.(float64)
	case tokenize.SLASH:
		if _, ok := leftVal.(string); ok {
			v.raiseError(errors.OperandMustBeNumber)
		}
		if _, ok := rightVal.(string); ok {
			v.raiseError(errors.OperandMustBeNumber)
		}
		return leftVal.(float64) / rightVal.(float64)
	case tokenize.MINUS:
		return leftVal.(float64) - rightVal.(float64)
	case tokenize.PLUS:
		if _, ok := rightVal.(string); ok {
			return fmt.Sprintf("%s%s", leftVal.(string), rightVal.(string))
		}
		return leftVal.(float64) + rightVal.(float64)
	case tokenize.GREATER:
		return leftVal.(float64) > rightVal.(float64)
	case tokenize.GREATER_EQUAL:
		return leftVal.(float64) >= rightVal.(float64)
	case tokenize.LESS_EQUAL:
		return leftVal.(float64) <= rightVal.(float64)
	case tokenize.LESS:
		return leftVal.(float64) < rightVal.(float64)
	case tokenize.BANG_EQUAL:
		return leftVal != rightVal
	case tokenize.EQUAL_EQUAL:
		return leftVal == rightVal
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
		if _, ok := strVal.(string); ok {
			v.raiseError(errors.OperandMustBeNumber)
		}
		return -expr.Right.Accept(v).(float64)
	}
}

func (v *Evaluator) raiseError(err string) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(70)
}
