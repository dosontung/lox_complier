package statement

import (
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/core"
)

type Interpreter struct {
	Evaluator core.ExprVisitor
}

func (i *Interpreter) VisitPrintStmt(statement *core.PrintStatement) interface{} {
	return i.valuate(statement.Expr)
}

func NewInterpreter(eval core.ExprVisitor) *Interpreter {
	return &Interpreter{Evaluator: eval}
}

func (i *Interpreter) valuate(statement core.Expression) interface{} {
	return statement.Accept(i.Evaluator)
}

func (i *Interpreter) Interpret(statement core.Statement) interface{} {
	return statement.Accept(i)
}

var _ core.StatementVisitor = (*Interpreter)(nil)
