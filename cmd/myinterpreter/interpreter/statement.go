package interpreter

import (
	"fmt"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/core"
)

var _ core.StatementVisitor = (*Interpreter)(nil)

func (v *Interpreter) VisitIfElseStmt(statement *core.IFElseStatement) {
	expr := v.Evaluate(statement.Expr)

	if expr.(bool) {
		v.Interpret(statement.ThenBranch)
	} else {
		if statement.ElseBranch != nil {
			v.Interpret(statement.ElseBranch)
		}

	}
}

func (v *Interpreter) VisitBlockStmt(statement *core.BlockStatement) {
	v.executeBlock(statement, NewEnvironment(v.env))
}

func (v *Interpreter) executeBlock(statement *core.BlockStatement, env *Environment) {
	previousEnv := v.env
	v.env = env
	for _, stmt := range statement.Statements {
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
