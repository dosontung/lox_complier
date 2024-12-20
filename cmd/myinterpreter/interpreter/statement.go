package interpreter

import (
	"fmt"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/core"
)

var _ core.StatementVisitor = (*Interpreter)(nil)

func (v *Interpreter) VisitFuncStmt(statement *core.FuncStatement) {
	env := v.env
	for env != nil {
		if err, _ := env.GetKey(statement.Name.Lexeme); err != nil {
			env.SetKey(statement.Name.Lexeme, statement)
		} else {
			v.raiseError("Duplicate name!")
		}
		env = env.Enclosing
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

	}
}

func (v *Interpreter) VisitWhileStmt(statement *core.WhileStatement) {
	for v.isTrue(v.Evaluate(statement.Expr)) {
		v.Interpret(statement.Body)
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
	v.executeBlock(statement.Statements, NewEnvironment(v.env))
}

func (v *Interpreter) executeBlock(statements []core.Statement, env *Environment) {
	previousEnv := v.env
	v.env = env
	for _, stmt := range statements {
		v.Interpret(stmt)
	}

	v.env = previousEnv
}

func (v *Interpreter) VisitVarDeclarationStmt(statement *core.VarDeclarationStatement) {
	//if statement.Expr != nil {
	value := v.Evaluate(statement.Expr)
	v.env.SetKey(statement.Name.Lexeme, value)
	return
	//}
	//if v.env.Enclosing == nil {
	//	v.env.SetKey(statement.Name.Lexeme, nil)
	//} else {
	//	v.env.Enclosing.SetKey(statement.Name.Lexeme, nil)
	//}

}

func (v *Interpreter) VisitExpressionStmt(statement *core.ExpressionStatement) {
	v.Evaluate(statement.Expr)
}

func (v *Interpreter) VisitPrintStmt(statement *core.PrintStatement) interface{} {
	value := v.Evaluate(statement.Expr)
	fmt.Println(value)
	return value
}

func (v *Interpreter) Interpret(statement core.Statement) interface{} {
	return statement.Accept(v)
}
