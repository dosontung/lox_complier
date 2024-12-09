package statement

import "errors"

type Environment struct {
	storage map[string]interface{}
}

func NewEnvironment() *Environment {
	return &Environment{make(map[string]interface{})}
}

func (env *Environment) SetKey(key string, value interface{}) {
	env.storage[key] = value
}

func (env *Environment) GetKey(key string) (error, interface{}) {
	if val, ok := env.storage[key]; ok {
		return nil, val
	}
	return errors.New("key not found"), nil
}
