package interpreter

type Interpreter struct {
	env *Environment
}

func NewInterpreter(env *Environment) *Interpreter {
	return &Interpreter{env: env}
}
