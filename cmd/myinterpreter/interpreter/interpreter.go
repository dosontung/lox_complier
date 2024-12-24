package interpreter

import (
	"errors"
)

type Interpreter struct {
	env *Environment
}

func NewInterpreter(env *Environment) *Interpreter {
	return &Interpreter{env: env}
}

func (v *Interpreter) SetKey(key string, val interface{}, assign bool) (interface{}, error) {
	env := v.env
	if assign == false {
		env.SetKey(key, val)
		return val, nil
	}
	for env != nil {
		if _, err := env.GetKey(key); err == nil {
			env.SetKey(key, val)
			return val, nil
		}
		env = env.Enclosing
	}
	return val, nil
}

func (v *Interpreter) GetKey(key string) (interface{}, error) {
	env := v.env
	for env != nil {
		if value, err := env.GetKey(key); err == nil {
			return value, nil
		}
		env = env.Enclosing
	}
	return nil, errors.New("key not found")
}
