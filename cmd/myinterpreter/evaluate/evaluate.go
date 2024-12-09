package evaluate

import (
	"fmt"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/core"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/errors"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/statement"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/tokenize"
	"math"
	"os"
	"reflect"
)

type Evaluator struct {
	env *statement.Environment
}

func NewEvaluator(env *statement.Environment) *Evaluator {
	return &Evaluator{env: env}
}

var _ core.ExprVisitor = &Evaluator{}

func (v *Evaluator) VisitVarExpr(expr *core.VarExpression) interface{} {
	err, i := v.env.GetKey(expr.Name.Lexeme)
	if err != nil {
		v.raiseError(errors.UndefinedVar, fmt.Sprintf(" '%s'.", expr.Name.Lexeme))
	}
	return i
}

func (v *Evaluator) VisitBinaryExpr(expr *core.BinaryExpression) interface{} {
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

func (v *Evaluator) VisitGroupingExpr(expr *core.GroupExpression) interface{} {
	return expr.Expr.Accept(v)
}

func (v *Evaluator) VisitLiteralExpr(expr *core.LiteralExpression) interface{} {
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

func (v *Evaluator) VisitUnaryExpr(expr *core.UnaryExpression) interface{} {
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

func (v *Evaluator) Evaluate(expr core.Expression) interface{} {
	return expr.Accept(v)
}

func (v *Evaluator) raiseError(err errors.CError, etcs ...string) {
	fmt.Fprint(os.Stderr, err)
	for _, etc := range etcs {
		fmt.Fprint(os.Stderr, etc)
	}
	fmt.Print("\n")
	os.Exit(70)
}
