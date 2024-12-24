package interpreter

import (
	"errors"
	"fmt"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/core"
)

var _ core.StatementVisitor = (*Interpreter)(nil)

func (v *Interpreter) VisitReturnStmt(statement *core.ReturnStatement) {
	if statement.Expr == nil {
		v.SetKey("Return", nil, false)
		return
	}
	retunVal := v.Evaluate(statement.Expr)
	v.SetKey("Return", retunVal, false)
}

func (v *Interpreter) VisitFuncStmt(statement *core.FuncStatement) {
	if _, err := v.GetKey(statement.Name.Lexeme); err != nil {
		funClone := statement.Clone()
		funClone.SetLocalVar(v.env)
		v.SetKey(statement.Name.Lexeme, funClone, false)
		//if statement.Name.Lexeme == "filter" {
		//	fmt.Printf("filter env %p\n", v.env)
		//	fmt.Printf("filter env %v\n", v.env)
		//}
	} else {
		v.raiseError("Duplicate name!")
	}

}

func (v *Interpreter) VisitForStmt(statement *core.ForStatement) {
	if statement.VarStatment != nil {
		v.Interpret(statement.VarStatment)
	}
	for v.isTrue(v.Evaluate(statement.Expr1)) {
		v.Interpret(statement.Body)
		if statement.Expr2 != nil {
			v.Evaluate(statement.Expr2)
		}
		// Break when get return in func
		if _, err := v.GetKey("Return"); err == nil {
			return
		}
	}
}

func (v *Interpreter) VisitWhileStmt(statement *core.WhileStatement) {
	for v.isTrue(v.Evaluate(statement.Expr)) {
		v.Interpret(statement.Body)
		// Break when get return in func
		if _, err := v.GetKey("Return"); err == nil {
			return
		}

	}
}

func (v *Interpreter) VisitIfElseStmt(statement *core.IFElseStatement) {
	expr := v.Evaluate(statement.Expr)
	if lvl, ok := expr.(bool); ok && lvl {
		v.Interpret(statement.ThenBranch)
		return
	} else if lvl, ok := expr.(float64); ok && lvl != 0 {
		v.Interpret(statement.ThenBranch)
		return
	} else if _, ok := expr.(string); ok {
		v.Interpret(statement.ThenBranch)
		return
	}
	if statement.ElseBranch != nil {
		v.Interpret(statement.ElseBranch)
	}

}

func (v *Interpreter) VisitBlockStmt(statement *core.BlockStatement) {
	if ret, ok := v.executeBlock(statement.Statements, NewEnvironment(v.env)); ok != nil {
		v.SetKey("Return", ret, false)
	}
}

func (v *Interpreter) executeBlock(statements []core.Statement, env *Environment) (interface{}, error) {
	previousEnv := v.env
	v.env = env
	defer func() { v.env = previousEnv }()
	for _, stmt := range statements {
		v.Interpret(stmt)
		if pushReturn, err := v.GetKey("Return"); err == nil {
			//previousEnv.SetKey("Return", pushReturn)
			return pushReturn, errors.New("return statement detected")
		}
	}
	//v.env = previousEnv
	return nil, nil
}

func (v *Interpreter) VisitVarDeclarationStmt(statement *core.VarDeclarationStatement) {
	value := v.Evaluate(statement.Expr)
	v.SetKey(statement.Name.Lexeme, value, false)
	return

}

func (v *Interpreter) VisitExpressionStmt(statement *core.ExpressionStatement) {
	v.Evaluate(statement.Expr)
}

func (v *Interpreter) VisitPrintStmt(statement *core.PrintStatement) interface{} {
	value := v.Evaluate(statement.Expr)
	if fun, ok := value.(*core.FuncStatement); ok {
		fmt.Printf("<fn %s>\n", fun.Name.Lexeme)
	} else {
		fmt.Println(value)
	}
	return value
}

func (v *Interpreter) Interpret(statement core.Statement) interface{} {
	return statement.Accept(v)
}
