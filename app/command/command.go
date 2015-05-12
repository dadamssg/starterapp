package command

import (
	"errors"
)

type Errorable interface {
	AddError(err CommandError)
	HasErrors() bool
	GetErrors() []CommandError
}

type CommandError struct {
	Code int
	Err  error
}

type Command struct {
	errors []CommandError
}

func (command *Command) AddError(err CommandError) {
	command.errors = append(command.errors, err)
}

func (command *Command) HasErrors() bool {
	return len(command.errors) > 0
}

func (command *Command) GetErrors() []CommandError {
	return command.errors
}

func AddCommandError(command Errorable, code int, err string) {
	command.AddError(CommandError{Code: code, Err: errors.New(err)})
}
