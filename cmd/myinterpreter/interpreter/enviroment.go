package interpreter

import "errors"

type Environment struct {
	storage   map[string]interface{}
	Enclosing *Environment
}

func NewEnvironment(env *Environment) *Environment {
	return &Environment{make(map[string]interface{}), env}
}

func (env *Environment) SetKey(key string, value interface{}) {
	env.storage[key] = value
}

func (env *Environment) GetKey(key string) (error, interface{}) {
	if val, ok := env.storage[key]; ok {
		if val == nil {
			val = "nil"
		}
		return nil, val
	}
	return errors.New("key not found"), nil
}
