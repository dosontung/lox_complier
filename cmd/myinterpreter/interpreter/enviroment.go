package interpreter

import (
	"errors"
)

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

func (env *Environment) GetKey(key string) (interface{}, error) {
	if val, ok := env.storage[key]; ok {
		if val == nil {
			val = "nil"
		}
		return val, nil
	}
	return nil, errors.New("key not found")
}
