package command

import (
	"github.com/dadamssg/commandbus"
	"reflect"
	"sync"
)

type ValidatorFunc func(cmd interface{}) []CommandError

type Validator struct {
	validatorFuncs map[reflect.Type]ValidatorFunc
	lock           sync.Mutex
}

func (v *Validator) Register(cmd interface{}, f ValidatorFunc) {
	v.lock.Lock()
	defer v.lock.Unlock()

	v.validatorFuncs[reflect.TypeOf(cmd)] = f
}

func NewValidator() *Validator {
	return &Validator{
		validatorFuncs: make(map[reflect.Type]ValidatorFunc),
	}
}

func RegisterMiddleware(v *Validator) commandbus.MiddlewareFunc {
	return func(cmd interface{}, next commandbus.HandlerFunc) {
		command, ok := cmd.(Errorable)

		if !ok {
			next(cmd)
			return
		}

		t := reflect.TypeOf(cmd)

		if fn, ok := v.validatorFuncs[t]; ok {
			if errs := fn(cmd); len(errs) > 0 {
				for _, err := range errs {
					command.AddError(err)
				}
				return
			}
		}

		next(cmd)
	}
}
