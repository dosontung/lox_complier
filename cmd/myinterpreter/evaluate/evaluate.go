package evaluate

import (
	"fmt"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/errors"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/parser"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/tokenize"
	"math"
	"os"
	"reflect"
)

type Evaluator struct {
}

func (v *Evaluator) VisitBinaryExpr(expr *parser.BinaryExpression) interface{} {
	rightVal := expr.Right.Accept(v)
	leftVal := expr.Left.Accept(v)
	sameType, isNumber := false, false
	if reflect.TypeOf(leftVal) == reflect.TypeOf(rightVal) {
		sameType = true
	}
	if _, ok := leftVal.(float64); ok {
		isNumber = true
	}
	switch expr.Operator.Type {
	case tokenize.STAR:
		if !isNumber || !sameType {
			v.raiseError(errors.OperandMustBeNumber)
		}
		return rightVal.(float64) * leftVal.(float64)
	case tokenize.SLASH:
		if !isNumber || !sameType {
			v.raiseError(errors.OperandMustBeNumber)
		}
		return leftVal.(float64) / rightVal.(float64)
	case tokenize.MINUS:
		if !sameType {
			v.raiseError(errors.OperandsMustBeSameType)
		}
		return leftVal.(float64) - rightVal.(float64)
	case tokenize.PLUS:
		if !sameType {
			v.raiseError(errors.OperandsMustBeSameType)
		}
		if _, ok := rightVal.(string); ok {
			return fmt.Sprintf("%s%s", leftVal.(string), rightVal.(string))
		}
		return leftVal.(float64) + rightVal.(float64)
	case tokenize.GREATER:
		if !isNumber || !sameType {
			v.raiseError(errors.OperandMustBeNumber)
		}
		return leftVal.(float64) > rightVal.(float64)
	case tokenize.GREATER_EQUAL:
		if !isNumber || !sameType {
			v.raiseError(errors.OperandMustBeNumber)
		}
		return leftVal.(float64) >= rightVal.(float64)
	case tokenize.LESS_EQUAL:
		if !isNumber || !sameType {
			v.raiseError(errors.OperandMustBeNumber)
		}
		return leftVal.(float64) <= rightVal.(float64)
	case tokenize.LESS:
		if !isNumber || !sameType {
			v.raiseError(errors.OperandMustBeNumber)
		}
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
		if _, ok := strVal.(float64); !ok {
			v.raiseError(errors.OperandMustBeNumber)
		}
		return -expr.Right.Accept(v).(float64)
	}
}

func (v *Evaluator) raiseError(err errors.CError) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(70)
}
