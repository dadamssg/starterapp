package user

import (
	"github.com/dadamssg/commandbus"
	"github.com/dadamssg/starterapp/app/command"
)

func Connect(bus *commandbus.CommandBus, validator *command.Validator, users UserRepository, mailer UserMailer) {
	registerUserHandler(bus, users, mailer)
	findUserByIdHandler(bus, users)
	registerUserValidator(validator, users)
}
