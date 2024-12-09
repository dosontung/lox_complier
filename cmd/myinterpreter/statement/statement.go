package statement

import (
	"fmt"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/core"
)

type Interpreter struct {
	Evaluator core.ExprVisitor
	env       *Environment
}

func (i *Interpreter) VisitVarDeclarationStmt(statement *core.VarDeclarationStatement) {
	if statement.Expr != nil {
		value := i.valuate(statement.Expr)
		i.env.SetKey(statement.Name.Lexeme, value)
	}
}

func (i *Interpreter) VisitExpressionStmt(statement *core.ExpressionStatement) {
	i.valuate(statement.Expr)
}

func (i *Interpreter) VisitPrintStmt(statement *core.PrintStatement) interface{} {
	value := i.valuate(statement.Expr)
	fmt.Println(value)
	return value
}

func NewInterpreter(eval core.ExprVisitor, env *Environment) *Interpreter {
	return &Interpreter{Evaluator: eval, env: env}
}

func (i *Interpreter) valuate(statement core.Expression) interface{} {
	return statement.Accept(i.Evaluator)
}

func (i *Interpreter) Interpret(statement core.Statement) interface{} {
	return statement.Accept(i)
}

var _ core.StatementVisitor = (*Interpreter)(nil)
