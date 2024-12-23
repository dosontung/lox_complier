package interpreter

import (
	"fmt"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/core"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/errors"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/tokenize"
	"math"
	"os"
	"reflect"
	"time"
)

var _ core.ExprVisitor = &Interpreter{}

func (v *Interpreter) raiseError(err errors.CError, etcs ...string) {
	fmt.Fprint(os.Stderr, err)
	for _, etc := range etcs {
		fmt.Fprint(os.Stderr, etc)
	}
	fmt.Print("\n")
	os.Exit(70)
}
func (v *Interpreter) VisitCallExpr(expr *core.CallExpression) interface{} {
	callee := expr.Callee
	if val, ok := callee.(*core.VarExpression); ok {
		if val.Name.Lexeme == "clock" {
			return v.nativeCall()
		}
		if fun, err := v.GetKey(val.Name.Lexeme); err != nil {
			v.raiseError("No func")
		} else {
			if len(expr.Params) != len(fun.(*core.FuncStatement).Params) {
				v.raiseError("Wrong number of parameters")
			}
			return v.funCall(fun, expr.Params)
		}

	}
	return 0
}

func (v *Interpreter) nativeCall() interface{} {
	now := time.Now() // current local time
	sec := now.Unix() // number of seconds since January 1, 1970 UTC
	return float64(sec)
}

func (v *Interpreter) funCall(fun interface{}, params []core.Expression) interface{} {
	fun_ := fun.(*core.FuncStatement)
	funEnv := NewEnvironment(v.env)
	for i, param := range params {
		funEnv.SetKey(fun_.Params[i].Lexeme, v.Evaluate(param))
	}
	v.executeBlock(fun_.Body, funEnv)
	if val, err := funEnv.GetKey("Return"); err == nil {
		return val
	}

	return "nil"
}

func (v *Interpreter) VisitLogicalExpr(expr *core.LogicalExpression) interface{} {
	vl := v.Evaluate(expr.Left)
	op := expr.Operator.Type
	if op == tokenize.OR && v.isTrue(vl) {
		return vl
	} else if op == tokenize.AND && !v.isTrue(vl) {
		return vl
	} else {
		return v.Evaluate(expr.Right)
	}
}

func (v *Interpreter) isTrue(vl interface{}) bool {
	if lvl, ok := vl.(bool); ok && lvl == true {
		return true
	}
	if lvl, ok := vl.(float64); ok && lvl != 0 {
		return true
	}
	if _, ok := vl.(string); ok {
		return true
	}
	return false
}
func (v *Interpreter) VisitAssignExpr(expr *core.AssignExpression) interface{} {
	value := v.Evaluate(expr.Expr)
	if _, err := v.SetKey(expr.Name.Lexeme, value, true); err == nil {
		return value
	}
	v.raiseError("No variable!")
	return value
}

func (v *Interpreter) VisitVarExpr(expr *core.VarExpression) interface{} {
	var i interface{}
	if value, err := v.GetKey(expr.Name.Lexeme); err == nil {
		if fun, ok := value.(*core.FuncStatement); ok {
			return fmt.Sprintf("<fn %s>", fun.Name.Lexeme)
		}
		return value
	}

	v.raiseError(errors.UndefinedVar, fmt.Sprintf(" '%s'.", expr.Name.Lexeme))
	return i
}

func (v *Interpreter) VisitBinaryExpr(expr *core.BinaryExpression) interface{} {
	rightVal := v.Evaluate(expr.Right)
	leftVal := v.Evaluate(expr.Left)
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

func (v *Interpreter) VisitGroupingExpr(expr *core.GroupExpression) interface{} {
	return expr.Expr.Accept(v)
}

func (v *Interpreter) VisitLiteralExpr(expr *core.LiteralExpression) interface{} {
	strVal := expr.Value
	switch strVal.(type) {
	//case string:
	//	return strVal
	case float64:
		if strVal == math.Trunc(strVal.(float64)) {
			return math.Trunc(strVal.(float64))
		}
		return strVal
	}
	return strVal
}

func (v *Interpreter) VisitUnaryExpr(expr *core.UnaryExpression) interface{} {
	strVal := expr.Right.Accept(v)
	switch expr.Operator.Type {
	case tokenize.BANG:
		return !v.isTrue(strVal)
	default: // tokenize.MINUS
		if _, ok := strVal.(float64); !ok {
			v.raiseError(errors.OperandMustBeNumber)
		}
		return -expr.Right.Accept(v).(float64)
	}
}

func (v *Interpreter) Evaluate(expr core.Expression) interface{} {
	return expr.Accept(v)
}
