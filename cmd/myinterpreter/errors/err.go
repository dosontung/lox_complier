package errors

type CError string

const (
	OperandMustBeNumber    CError = "Operand must be a number."
	OperandsMustBeSameType CError = "Operands must be two numbers or two strings."
	UndefinedVar           CError = "Undefined variable"
)
